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
	loggingService      *logging.ScrapeLoggingService
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
	log.Printf("SMQ HOST: %s:%s ", smqHost, smqPort)
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

func (service SessionJobHandler) GetAdsPageJobResult() chan jobs.AdsPageJobResult {
	return service.resultsChannel
}

func (service SessionJobHandler) GetResultsChannel() chan icollector.AdsResults {
	return service.adsResultsChannel
}

func (service SessionJobHandler) readResultsFromServices() {
	for _, scrapingService := range service.marketServiceMapper.GetAllServices() {
		ss := scrapingService
		//wg.Add(1)
		go func() {
			for {
				adsResult := <-*ss.GetResultsChannel()
				service.resultsChannel <- adsResult
			}

		}()
	}
}

func (service SessionJobHandler) StartWithoutMQ() {
	log.Println("Session Job Handler Starting without MQ")
	//service.messageQueueService.Start()
	// read results from services
	service.assignJobsToServices()
	service.readResultsFromServices()

	//wg.Add(1)
	go func() {
		//defer wg.Done()
		for {
			select {
			case job := <-service.jobChannel:
				//service.AddScrapingJobToScrapingService(job)
				service.AddScrapingJob(job)
			case res := <-service.resultsChannel:
				go func() {
					//service.processResultsNOPublish(res)
					//service.processResults(res)
					service.processResultsAndAddInternal(res)
				}()
			case <-service.context.Done():
				log.Println("Session Job Handler Terminating...")
				return
			}
		}
	}()
}

func (service SessionJobHandler) Start() {
	log.Println("Session Job Handler Starting ")
	service.messageQueueService.Start()
	service.assignJobsToServices()
	for _, scrapingService := range service.marketServiceMapper.GetAllServices() {
		ss := scrapingService
		//wg.Add(1)
		go func() {
			for {
				tmp := ss.GetResultsChannel()
				adsResult := <-*tmp
				service.resultsChannel <- adsResult
			}

		}()
	}
	go func() {
		for {
			select {
			case job := <-service.jobChannel:
				//service.AddScrapingJobToScrapingService(job)
				service.AddScrapingJob(job)
			case res := <-service.resultsChannel:
				go func() {
					service.processResults(res)
				}()
			case <-service.context.Done():
				log.Println("Session Job Handler Terminating...")
				return
			}
		}
	}()

	go func() {
		for {
			service.getJobFromMQ()
		}
	}()
}

func (service SessionJobHandler) Start_old() {
	log.Println("Session Job Handler Service Start")

	log.Println("start waiting for signal")

	service.messageQueueService.Start()

	go func() {
		for {
			service.getJobFromMQ()
		}
	}()

	go func() {
		for {
			//service.processJob()
		}
	}()

	go func() {
		for {
			//service.sendResults()
		}
	}()

	go func() {
		for {
			select {
			case <-service.context.Done():
				log.Println("Session Job Handler Terminating...")
				return

			}
		}
	}()
}

func (service SessionJobHandler) getJobFromMQ() {
	// pop message from MQ
	message := service.messageQueue.GetMessageWithDelete(service.jobsTopicName)
	var scrapeJob jobs.SessionJob
	if len(*message) > 0 {
		log.Println("GOT MESSAGE FROM MQ ")
		err := json.Unmarshal(*message, &scrapeJob)
		log.Println("JOB FROM MQ: ", scrapeJob.ToString())
		if err != nil {
			// push message back in the queue
			service.messageQueue.PutMessage(service.jobsTopicName, *message)
			panic(err)
		}
		service.jobChannel <- scrapeJob
	}

}

// AddScrapingJob adds the job to the session handler jobs channel
func (service SessionJobHandler) AddScrapingJob(job jobs.SessionJob) {
	service.jobChannel <- job
}

func (service SessionJobHandler) assignJobsToServices() {
	go func() {
		for {
			job := <-service.jobChannel
			scrapingService := service.marketServiceMapper.GetScrapingService(job.Market.Name)
			go func() {
				//log.Printf("ASSIGN JOB TO SCRAPING SERVICE : %s ", job.Market.Name)
				scrapingService.AddJob(job)
			}()
		}
	}()
}

//func (service SessionJobHandler) AddScrapingJobToScrapingService(job jobs.SessionJob) {
//	log.Println("Get scraping service for market : ", job.Market.Name)
//	s := service.marketServiceMapper.GetScrapingService(job.Market.Name)
//	s.AddJob(job)
//}

func (service SessionJobHandler) processResultsNOPublish(result jobs.AdsPageJobResult) {
	var foundAds []jobs.Ad
	for _, res := range *result.Data {
		fAd := newFromAd(res, result.RequestedScrapingJob.Criteria.Brand, result.RequestedScrapingJob.Criteria.CarModel, result.RequestedScrapingJob.Criteria.Fuel)
		foundAds = append(foundAds, fAd)
	}

	result.Data = &foundAds

	// send result to MQ
	if !result.IsLastPage {
		newJob := result.RequestedScrapingJob
		newJob.Market.PageNumber++
		//service.AddScrapingJobToScrapingService(newJob)
		service.AddScrapingJob(newJob)
	} else {
		log.Println("Got LAST PAGE")
	}

	//log.Println("Job handler done processing: ", len(*result.Data), " is last page ", result.IsLastPage)

}

func (service SessionJobHandler) processResults(result jobs.AdsPageJobResult) {
	var foundAds []jobs.Ad
	if result.Data != nil {
		for _, res := range *result.Data {
			fAd := newFromAd(res, result.RequestedScrapingJob.Criteria.Brand, result.RequestedScrapingJob.Criteria.CarModel, result.RequestedScrapingJob.Criteria.Fuel)
			foundAds = append(foundAds, fAd)
		}

		result.Data = &foundAds
	}

	// send result to MQ
	if !result.IsLastPage {
		newJob := service.createNewJobFromResult(result)
		//service.AddScrapingJobToScrapingService(newJob)
		service.pushSessionJobToMQ(newJob)
	} else {
		log.Println("Got LAST PAGE")
	}
	//log.Println("Sending to queue ", len(*result.Data), " is last page ", result.IsLastPage, " JOB: ", result.RequestedScrapingJob.ToString())

	service.messageQueueService.PublishResults(result)
}

func (service SessionJobHandler) processResultsAndAddInternal(result jobs.AdsPageJobResult) {
	var foundAds []jobs.Ad
	if result.Data == nil {
		fakeData := []jobs.Ad{}
		result.Data = &fakeData
	}
	for _, res := range *result.Data {
		fAd := newFromAd(res, result.RequestedScrapingJob.Criteria.Brand, result.RequestedScrapingJob.Criteria.CarModel, result.RequestedScrapingJob.Criteria.Fuel)
		foundAds = append(foundAds, fAd)
	}

	result.Data = &foundAds

	// send result to MQ
	if !result.IsLastPage {
		newJob := service.createNewJobFromResult(result)
		//service.AddScrapingJobToScrapingService(newJob)
		service.AddScrapingJob(newJob)
	} else {
		log.Println("Got LAST PAGE")
	}
	log.Println("Sending to queue ", len(*result.Data), " is last page ", result.IsLastPage, " JOB ", result.RequestedScrapingJob.ToString())

	service.messageQueueService.PublishResults(result)
}

func (service SessionJobHandler) createNewJobFromResult(result jobs.AdsPageJobResult) jobs.SessionJob {
	newJob := result.RequestedScrapingJob
	newJob.Market.PageNumber++
	return newJob
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

func (service SessionJobHandler) listResults(res jobs.AdsPageJobResult) {

	//res := <-service.resultsChannel
	log.Println("Showing page : ", res.PageNumber)
	log.Println("Is Last page : ", res.IsLastPage)
	//for _, ad := range *res.Data {
	//	log.Printf("%+v", ad)
	//}
}

func (service SessionJobHandler) pushSessionJobToMQ(job jobs.SessionJob) {
	jobBytes, err := json.Marshal(&job)
	if err != nil {
		panic(err)
	}
	log.Printf("SESSION HANDLER : PUSHING JOB %s", job.ToString())
	service.messageQueue.PutMessage(service.jobsTopicName, jobBytes)
}

func (service SessionJobHandler) AddJobToMQ(job jobs.SessionJob) {
	service.pushSessionJobToMQ(job)
}
