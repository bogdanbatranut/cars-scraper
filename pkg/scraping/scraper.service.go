package scraping

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/repos"
	"carscraper/pkg/scraping/strategies"
	"carscraper/pkg/scraping/urlbuilder"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type PageScrapingService struct {
	messageQueue         repos.IMessageQueue
	requestedScrapingJob jobs.SessionJob
	jobChannel           chan jobs.SessionJob
	//messageChannel       chan []byte
	resultsChannel         chan jobs.AdsPageJobResult
	pagesToScrapeTopicName string
	resultsTopicName       string
}

type PageScrapingServiceConfiguration func(sjc *PageScrapingService)

func NewPageScrapingService(cfgs ...PageScrapingServiceConfiguration) *PageScrapingService {
	service := &PageScrapingService{
		jobChannel:     make(chan jobs.SessionJob),
		resultsChannel: make(chan jobs.AdsPageJobResult),
		//messageChannel: make(chan []byte),
		pagesToScrapeTopicName: "pagesToScrape",
		resultsTopicName:       "results",
	}
	for _, cfg := range cfgs {
		cfg(service)
	}
	return service
}

func WithSimpleMessageQueueRepository1() PageScrapingServiceConfiguration {
	smqr := repos.NewSimpleMessageQueueRepository(
		"http://127.0.0.1:3333",
	)
	return WithMessageQueueRepository1(smqr)
}

func WithMessageQueueRepository1(mqr repos.IMessageQueue) PageScrapingServiceConfiguration {
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
		PageURL: urlbuilder.PageURL{
			MarketName: "",
			MarketURL:  "",
			CarBrand:   "",
			CarModel:   "",
			YearFrom:   nil,
			YearTo:     nil,
			Fuel:       nil,
			KmFrom:     nil,
			KmTo:       nil,
			PageNumber: 0,
		},
		Success: true,
		Data:    &pageResults,
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
	//return pageResults
}

func (sjc PageScrapingService) sendResults() {
	// all fine til here so push the results
	jobResult := <-sjc.resultsChannel
	resBytes, err := json.Marshal(&jobResult)
	if err != nil {
		panic(err)
	}
	sjc.messageQueue.PutMessage(sjc.resultsTopicName, resBytes)

	// add next page to pagesToScrape
	if jobResult.PageURL.PageNumber > 0 {

	}

	//sjc.messageQueue.PutMessage()

}
