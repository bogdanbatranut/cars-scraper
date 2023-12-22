package scraper

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/jobs"
	"carscraper/pkg/repos"
	"carscraper/pkg/scraping/markets"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
)

type PageScrapingService struct {
	messageQueue          repos.IMessageQueue
	requestedScrapingJob  jobs.SessionJob
	jobChannel            chan jobs.SessionJob
	additionalJobsChannel chan jobs.SessionJob
	//messageChannel       chan []byte
	resultsChannel         chan jobs.AdsPageJobResult
	pagesToScrapeTopicName string
	resultsTopicName       string
	jobsTopicName          string
}

type PageScrapingServiceConfiguration func(sjc *PageScrapingService)

func NewPageScrapingService(cfg amconfig.IConfig, cfgs ...PageScrapingServiceConfiguration) *PageScrapingService {
	jobsTopicName := cfg.GetString(amconfig.SMQJobsTopicName)
	service := &PageScrapingService{
		jobChannel:            make(chan jobs.SessionJob),
		resultsChannel:        make(chan jobs.AdsPageJobResult),
		additionalJobsChannel: make(chan jobs.SessionJob),
		//messageChannel: make(chan []byte),
		pagesToScrapeTopicName: jobsTopicName,
		resultsTopicName:       cfg.GetString(amconfig.SMQResultsTopicName),
		jobsTopicName:          cfg.GetString(amconfig.SMQJobsTopicName),
	}
	for _, cfg := range cfgs {
		cfg(service)
	}
	return service
}

func WithSimpleMessageQueueRepository(cfg amconfig.IConfig) PageScrapingServiceConfiguration {
	smqHost := cfg.GetString(amconfig.SMQURL)
	smqPort := cfg.GetString(amconfig.SMQHTTPPort)
	smqr := repos.NewSimpleMessageQueueRepository(fmt.Sprintf("http://%s:%s", smqHost, smqPort))
	return WithMessageQueueRepository(smqr)
}

func WithMessageQueueRepository(mqr repos.IMessageQueue) PageScrapingServiceConfiguration {
	return func(cis *PageScrapingService) {
		cis.messageQueue = mqr
	}
}

func (sjc PageScrapingService) Start() {
	log.Println("Scraping Service Start")
	done := make(chan bool, 1)

	signalsChannel := make(chan os.Signal, 1)
	signal.Notify(signalsChannel, syscall.SIGINT, syscall.SIGTERM)
	log.Println("start waiting for signal")

	_, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			sjc.getJobFromMQ()
		}
	}()

	go func() {
		for {
			sjc.processJob()
		}
	}()

	go func() {
		for {
			sjc.sendResults()
		}
	}()

	go func() {
		for {
			sjc.pushAdditionalSessionJob()
		}
	}()

	go func() {
		log.Println("Waiting for signal")
		sig := <-signalsChannel
		log.Println("Got signal:", sig)
		log.Println("Terminating...")
		cancel()
		done <- true
	}()

	<-done

}

func (sjc PageScrapingService) getJobFromMQ() {
	// pop message from MQ
	message := sjc.messageQueue.GetMessageWithDelete(sjc.pagesToScrapeTopicName)
	var scrapeJob jobs.SessionJob
	if len(*message) > 0 {
		err := json.Unmarshal(*message, &scrapeJob)
		if err != nil {
			// push message back in the queue
			sjc.messageQueue.PutMessage(sjc.pagesToScrapeTopicName, *message)
			panic(err)
		}
		//log.Println("pushing job to jobChannel")
		//log.Printf("Scraping service GOT job: criteria: %d, market: %d, pageNumber: %d", scrapeJob.CriteriaID, scrapeJob.MarketID, scrapeJob.Market.PageNumber)

		sjc.jobChannel <- scrapeJob
	}
}

func (sjc PageScrapingService) processJob() {
	// crawl the page
	job := <-sjc.jobChannel
	marketName := job.Market.Name

	availableImplementations := markets.NewImplemetationStrategies()
	implementation := availableImplementations.GetImplementation(marketName)

	if implementation != nil {
		pageResults, isLastPage, err := implementation.Execute(job)

		log.Printf("Total Results for page: %d isLastpage : %t", len(pageResults), isLastPage)

		if pageResults == nil {
			log.Println("THERE ARE NO RESULTS SO RETURN...")
			return
		}

		jobResult := jobs.AdsPageJobResult{
			RequestedScrapingJob: job,
			IsLastPage:           isLastPage,
			Success:              true,
			Data:                 &pageResults,
			PageNumber:           job.Market.PageNumber,
		}

		// TODO if we have an error while scraping we need to see what happens...

		if err != nil {
			// push message back in the queue
			//message := <-sjc.messageChannel
			//sjc.messageQueue.PutMessage("requestedJobs", message)
			panic(err)
		}
		if jobResult.Data == nil {
			log.Printf("NO RESULTS !!!! %+v", jobResult)
		}
		sjc.resultsChannel <- jobResult

		// determine here if a new scrapejob should be created and create it
		if jobResult.IsLastPage {
			return
		}
		sjc.createNewSessionJob(job)
	}

}

func (sjc PageScrapingService) createNewSessionJob(oldJob jobs.SessionJob) {
	pageNumber := oldJob.Market.PageNumber
	pageNumber++

	newMarket := jobs.Market{
		Name:       oldJob.Market.Name,
		PageNumber: pageNumber,
	}

	additionalJob := jobs.SessionJob{
		SessionID:  oldJob.SessionID,
		JobID:      uuid.New(),
		CriteriaID: oldJob.CriteriaID,
		MarketID:   oldJob.MarketID,
		Criteria:   oldJob.Criteria,
		Market:     newMarket,
	}
	sjc.additionalJobsChannel <- additionalJob
}

func (sjc PageScrapingService) sendResults() {
	// all fine til here so push the results
	jobResult := <-sjc.resultsChannel

	log.Printf("TotalResults SENT : %d: ", len(*jobResult.Data))
	//log.Printf("Scraping service RESULTS: criteria: %d, market: %d, pageNumber: %d", jobResult.RequestedScrapingJob.CriteriaID, jobResult.RequestedScrapingJob.MarketID, jobResult.RequestedScrapingJob.Market.PageNumber)

	resBytes, err := json.Marshal(&jobResult)
	if err != nil {
		panic(err)
	}
	sjc.messageQueue.PutMessage(sjc.resultsTopicName, resBytes)
	//sjc.messageQueue.PutMessage(jobResult.GetTopic(), resBytes)
}

func (sjc PageScrapingService) pushAdditionalSessionJob() {
	job := <-sjc.additionalJobsChannel
	//log.Printf("Scraping service ADDITIONAL: criteria: %d, market: %d, pageNumber: %d", job.CriteriaID, job.MarketID, job.Market.PageNumber)

	jobBytes, err := json.Marshal(&job)
	if err != nil {
		panic(err)
	}
	sjc.messageQueue.PutMessage(sjc.jobsTopicName, jobBytes)
}
