package scraper

import (
	"carscraper/pkg/config"
	"carscraper/pkg/jobs"
	"carscraper/pkg/repos"
	"carscraper/pkg/scraping/strategies"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

func NewPageScrapingService(cfg config.IConfig, cfgs ...PageScrapingServiceConfiguration) *PageScrapingService {
	jobsTopicName := cfg.GetString(config.SMQJobsTopicName)
	service := &PageScrapingService{
		jobChannel:            make(chan jobs.SessionJob),
		resultsChannel:        make(chan jobs.AdsPageJobResult),
		additionalJobsChannel: make(chan jobs.SessionJob),
		//messageChannel: make(chan []byte),
		pagesToScrapeTopicName: jobsTopicName,
		resultsTopicName:       cfg.GetString(config.SMQResultsTopicName),
		jobsTopicName:          cfg.GetString(config.SMQJobsTopicName),
	}
	for _, cfg := range cfgs {
		cfg(service)
	}
	return service
}

func WithSimpleMessageQueueRepository(cfg config.IConfig) PageScrapingServiceConfiguration {
	smqHost := cfg.GetString(config.SMQURL)
	smqPort := cfg.GetString(config.SMQHTTPPort)
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
			time.Sleep(2 * time.Second)
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
		log.Println("pushing job to jobChannel")
		log.Println(scrapeJob)
		sjc.jobChannel <- scrapeJob
		//sjc.messageChannel <- *message
	}

	//return scrapeJob, message
}

func (sjc PageScrapingService) processJob() {
	// crawl the page
	job := <-sjc.jobChannel

	marketName := job.Market.Name

	availableImplementations := strategies.NewImplemetationStrategies()
	implementation := availableImplementations.GetImplementation(marketName)

	pageResults, err := implementation.Execute(job.Market.URL)

	jobResult := jobs.AdsPageJobResult{
		RequestedScrapingJob: job,
		IsLastPage:           false,
		Success:              true,
		Data:                 &pageResults,
	}

	if job.PageNumber == 3 {
		jobResult.IsLastPage = true
	}

	if err != nil {
		// push message back in the queue
		//message := <-sjc.messageChannel
		//sjc.messageQueue.PutMessage("requestedJobs", message)
		panic(err)
	}
	log.Println("Sending results to channel")
	log.Println(">>>> ", pageResults)
	sjc.resultsChannel <- jobResult

	// determine here if a new scrapejob should be created and create it
	if jobResult.IsLastPage {
		return
	}
	sjc.createNewSessionJob(job)
	//return pageResults
}

func (sjc PageScrapingService) createNewSessionJob(oldJob jobs.SessionJob) {
	pageNuber := oldJob.PageNumber
	pageNuber++

	additionalJob := jobs.SessionJob{
		SessionID:  oldJob.SessionID,
		JobID:      uuid.New(),
		CriteriaID: oldJob.CriteriaID,
		MarketID:   oldJob.MarketID,
		Criteria:   oldJob.Criteria,
		Market:     oldJob.Market,
		PageNumber: pageNuber,
	}
	sjc.additionalJobsChannel <- additionalJob
}

func (sjc PageScrapingService) sendResults() {
	// all fine til here so push the results
	jobResult := <-sjc.resultsChannel
	resBytes, err := json.Marshal(&jobResult)
	if err != nil {
		panic(err)
	}
	sjc.messageQueue.PutMessage(sjc.resultsTopicName, resBytes)
}

func (sjc PageScrapingService) pushAdditionalSessionJob() {
	job := <-sjc.additionalJobsChannel
	jobBytes, err := json.Marshal(&job)
	if err != nil {
		panic(err)
	}
	sjc.messageQueue.PutMessage(sjc.jobsTopicName, jobBytes)
}
