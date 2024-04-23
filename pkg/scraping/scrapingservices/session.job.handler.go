package scrapingservices

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/mq"
	"carscraper/pkg/repos"
	"carscraper/pkg/scraping/icollector"
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// SessionJobHandler is a service that allows jobs to be added and will forward the jobs
// to the correct scraping service
// when the job is done it will forward the results to interested parties
type SessionJobHandler struct {
	context             context.Context
	messageQueue        repos.IMessageQueue
	jobChannel          chan jobs.SessionJob
	resultsChannel      chan jobs.AdsPageJobResult
	adsResultsChannel   chan icollector.AdsResults
	resultsTopicName    string
	jobsTopicName       string
	loggingService      logging.ScrapeLoggingService
	marketServiceMapper IScrapingServicesMapper
	messageQueueService *mq.MessageQueueService
	//scrapingServices    []IScrapingService
}

type SessionJobHandlerServiceConfiguration func(sjc *SessionJobHandler)

func WithMarketService(marketName string, service IScrapingService) SessionJobHandlerServiceConfiguration {
	return func(jobHandler *SessionJobHandler) {
		jobHandler.marketServiceMapper.AddMarketService(marketName, service)
	}
}

// func NewSessionJobHandler(ctx context.Context, cfg amconfig.IConfig, rodScraperService IScrapingService, collyScraperService IScrapingService, jsonScrapingService IScrapingService, cfgs ...SessionJobHandlerServiceConfiguration) *SessionJobHandler {
func NewSessionJobHandler(ctx context.Context, cfg amconfig.IConfig, cfgs ...SessionJobHandlerServiceConfiguration) *SessionJobHandler {
	smqHost := cfg.GetString(amconfig.SMQURL)
	smqPort := cfg.GetString(amconfig.SMQHTTPPort)
	smqr := repos.NewSimpleMessageQueueRepository(fmt.Sprintf("http://%s:%s", smqHost, smqPort))
	log.Println("Message Queue URL : ", fmt.Sprintf("http://%s:%s", smqHost, smqPort))

	marketServiceMapper := NewScrapingServicesMapper()

	//marketServiceMapper.AddMarketService("autoscout", rodScraperService)
	//marketServiceMapper.AddMarketService("mobile.de", collyScraperService)
	//marketServiceMapper.AddMarketService("autotracknl", rodScraperService)
	//marketServiceMapper.AddMarketService("autovit", jsonScrapingService)

	service := &SessionJobHandler{
		context:             ctx,
		resultsChannel:      make(chan jobs.AdsPageJobResult),
		jobChannel:          make(chan jobs.SessionJob),
		messageQueue:        smqr,
		resultsTopicName:    cfg.GetString(amconfig.SMQResultsTopicName),
		jobsTopicName:       cfg.GetString(amconfig.SMQJobsTopicName),
		loggingService:      logging.NewScrapeLoggingService(cfg),
		marketServiceMapper: marketServiceMapper,
		messageQueueService: mq.NewMessageQueueService(cfg, mq.WithProdMessageQueue()),
	}
	for _, cfg := range cfgs {
		cfg(service)
	}
	return service
}

func JobHandlerWithSimpleMessageQueueRepository(cfg amconfig.IConfig) SessionJobHandlerServiceConfiguration {
	smqHost := cfg.GetString(amconfig.SMQURL)
	smqPort := cfg.GetString(amconfig.SMQHTTPPort)
	smqr := repos.NewSimpleMessageQueueRepository(fmt.Sprintf("http://%s:%s", smqHost, smqPort))
	log.Println("Message Queue URL : ", fmt.Sprintf("http://%s:%s", smqHost, smqPort))
	return JobHandlerWithMessageQueueRepository(smqr)
}

func JobHandlerWithMessageQueueRepository(mqr repos.IMessageQueue) SessionJobHandlerServiceConfiguration {
	return func(cis *SessionJobHandler) {
		cis.messageQueue = mqr
	}
}

func (sjc SessionJobHandler) GetAdsPageJobResult() chan jobs.AdsPageJobResult {
	return sjc.resultsChannel
}

func (sjc SessionJobHandler) GetResultsChannel() chan icollector.AdsResults {
	return sjc.adsResultsChannel
}

func (sjc SessionJobHandler) StartWithoutMQ() {
	log.Println("Session Job Handler Starting without MQ")
	sjc.messageQueueService.Start()
	for _, scrapingService := range sjc.marketServiceMapper.GetAllServices() {
		ss := scrapingService
		//wg.Add(1)
		go func() {
			for {
				tmp := ss.GetResultsChannel()
				adsResult := <-*tmp
				sjc.resultsChannel <- adsResult
			}

		}()
	}

	//wg.Add(1)
	go func() {
		//defer wg.Done()
		for {
			select {
			case job := <-sjc.jobChannel:
				sjc.AddScrapingJob(job)
			case res := <-sjc.resultsChannel:
				go func() {
					sjc.processResultsNOPublish(res)
				}()
			case <-sjc.context.Done():
				log.Println("Session Job Handler Terminating...")
				return
			}
		}
	}()
}

func (sjc SessionJobHandler) Start() {
	log.Println("Session Job Handler Starting ")
	sjc.messageQueueService.Start()
	for _, scrapingService := range sjc.marketServiceMapper.GetAllServices() {
		ss := scrapingService
		//wg.Add(1)
		go func() {
			for {
				tmp := ss.GetResultsChannel()
				adsResult := <-*tmp
				sjc.resultsChannel <- adsResult
			}

		}()
	}

	//wg.Add(1)
	go func() {
		//defer wg.Done()
		for {
			select {
			case job := <-sjc.jobChannel:
				sjc.AddScrapingJob(job)
			case res := <-sjc.resultsChannel:
				go func() {
					sjc.processResults(res)
				}()
			case <-sjc.context.Done():
				log.Println("Session Job Handler Terminating...")
				return
			}
		}
	}()

	go func() {
		for {
			sjc.getJobFromMQ()
		}
	}()
	//wg.Wait()
}
func (sjc SessionJobHandler) Start_old() {
	log.Println("Session Job Handler Service Start")

	log.Println("start waiting for signal")

	sjc.messageQueueService.Start()

	go func() {
		for {
			sjc.getJobFromMQ()
		}
	}()

	go func() {
		for {
			//sjc.processJob()
		}
	}()

	go func() {
		for {
			//sjc.sendResults()
		}
	}()

	go func() {
		for {
			select {
			case <-sjc.context.Done():
				log.Println("Session Job Handler Terminating...")
				return

			}
		}
	}()
}

func (sjc SessionJobHandler) getJobFromMQ() {
	// pop message from MQ
	message := sjc.messageQueue.GetMessageWithDelete(sjc.jobsTopicName)
	var scrapeJob jobs.SessionJob
	if len(*message) > 0 {
		log.Printf("got message from the message queue")
		err := json.Unmarshal(*message, &scrapeJob)
		if err != nil {
			// push message back in the queue
			sjc.messageQueue.PutMessage(sjc.jobsTopicName, *message)
			panic(err)
		}
		sjc.jobChannel <- scrapeJob
	}
}

func (sjc SessionJobHandler) readResults() {
	for _, scrapingService := range sjc.marketServiceMapper.GetAllServices() {
		ss := scrapingService
		go func() {
			for {
				log.Println("Waiting for results ....")
				adsResult := <-*ss.GetResultsChannel()
				log.Println("Add results to session handler results channel")
				sjc.resultsChannel <- adsResult
			}
		}()

	}
	//log.Println("DONE READING ?????")
}

func (sjc SessionJobHandler) AddScrapingJob(job jobs.SessionJob) {
	log.Println("Get scraping service for market : ", job.Market.Name)
	s := sjc.marketServiceMapper.GetScrapingService(job.Market.Name)
	s.AddJob(job)
}

func (sjc SessionJobHandler) processResultsNOPublish(result jobs.AdsPageJobResult) {
	var foundAds []jobs.Ad
	for _, res := range *result.Data {
		fAd := newFromAd(res, result.RequestedScrapingJob.Criteria.Brand, result.RequestedScrapingJob.Criteria.CarModel, result.RequestedScrapingJob.Criteria.Fuel)
		foundAds = append(foundAds, fAd)
	}

	result.Data = &foundAds

	// send result to MQ
	log.Println("Session Job Handler : Processing results")
	if !result.IsLastPage {
		newJob := result.RequestedScrapingJob
		newJob.Market.PageNumber++
		sjc.AddScrapingJob(newJob)
	} else {
		log.Println("Got LAST PAGE")
	}

	for _, ad := range *result.Data {
		if ad.Title == nil {
			log.Println("TITLE NIL")
		} else {
			log.Printf("Ad Title: %s", *ad.Title)
		}
	}
	log.Println("Job handler done processing: ", len(*result.Data), " is last page ", result.IsLastPage)

}

func (sjc SessionJobHandler) processResults(result jobs.AdsPageJobResult) {
	var foundAds []jobs.Ad
	for _, res := range *result.Data {
		fAd := newFromAd(res, result.RequestedScrapingJob.Criteria.Brand, result.RequestedScrapingJob.Criteria.CarModel, result.RequestedScrapingJob.Criteria.Fuel)
		foundAds = append(foundAds, fAd)
	}

	result.Data = &foundAds

	// send result to MQ
	log.Println("Session Job Handler : Processing results")
	if !result.IsLastPage {
		newJob := result.RequestedScrapingJob
		newJob.Market.PageNumber++
		sjc.AddScrapingJob(newJob)
	} else {
		log.Println("Got LAST PAGE")
	}
	log.Println("Sending to queue ", len(*result.Data), " is last page ", result.IsLastPage)

	sjc.messageQueueService.PublishResults(result)
}

func newFromAd(ad jobs.Ad, brand string, model string, fuel string) jobs.Ad {
	newAd := jobs.Ad{
		Title:              ad.Title,
		Brand:              brand,
		Model:              model,
		Year:               ad.Year,
		Km:                 ad.Km,
		Fuel:               fuel,
		Price:              ad.Price,
		AdID:               ad.AdID,
		Ad_url:             ad.Ad_url,
		SellerType:         ad.SellerType,
		SellerName:         ad.SellerName,
		SellerNameInMarket: ad.SellerNameInMarket,
		SellerOwnURL:       ad.SellerOwnURL,
		SellerMarketURL:    ad.SellerMarketURL,
		Thumbnail:          ad.Thumbnail,
	}
	return newAd
}

func (sjc SessionJobHandler) listResults(res jobs.AdsPageJobResult) {

	//res := <-sjc.resultsChannel
	log.Println("Showing page : ", res.PageNumber)
	log.Println("Is Last page : ", res.IsLastPage)
	//for _, ad := range *res.Data {
	//	log.Printf("%+v", ad)
	//}
}
