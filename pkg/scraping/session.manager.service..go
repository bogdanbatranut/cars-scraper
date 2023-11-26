package scraping

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/repos"
	"carscraper/pkg/scraping/urlbuilder"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type SessionManagerService struct {
	messageQueue repos.IMessageQueue
	dps          *DonePagesService

	marketCriteriaTopicName string
	pagesToScrapeTopicName  string
	visitedPagesTopicName   string
}

type SessionManagerServiceConfiguration func(smsc *SessionManagerService)

func NewSessionManagerService(cfgs ...SessionManagerServiceConfiguration) *SessionManagerService {
	sms := &SessionManagerService{
		marketCriteriaTopicName: "criterias",
		pagesToScrapeTopicName:  "pagesToScrape",
		visitedPagesTopicName:   "scrapedPages",
		dps:                     NewDonePagesService(),
	}
	for _, cfg := range cfgs {
		cfg(sms)
	}
	return sms
}

func WithSimpleMessageQueueRepository2() SessionManagerServiceConfiguration {
	smqr := repos.NewSimpleMessageQueueRepository(
		"http://host.docker.internal:3333",
	)
	return func(sms *SessionManagerService) {
		sms.messageQueue = smqr
	}
}

func (sms SessionManagerService) Start() {
	log.Println("Session Management Service Start")
	done := make(chan bool, 1)

	signalsChannel := make(chan os.Signal, 1)
	signal.Notify(signalsChannel, syscall.SIGINT, syscall.SIGTERM)
	log.Println("start waiting for signal")

	_, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			sms.popMarketCriteria()
		}
	}()

	go func() {
		for {
			sms.checkCompletedPages()
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

// GetA Market Criteria and pushes the PageToScrape in the topic for pages to scrape
func (sms SessionManagerService) popMarketCriteria() {
	message := sms.messageQueue.GetMessageWithDelete(sms.marketCriteriaTopicName)
	var mc jobs.SessionJob
	if len(*message) > 0 {
		err := json.Unmarshal(*message, &mc)
		if err != nil {
			// push message back in the queue
			//sjc.messageQueue.PutMessage(sjc.pagesToScrapeTopicName, *message)
			panic(err)
		}
		// convert marketCriteria to page to scrape
		log.Println("pushing job to jobChannel")
		pageToScrape := urlbuilder.PageURL{
			MarketName: mc.Market.Name,
			MarketURL:  mc.Market.URL,
			CarBrand:   mc.Criteria.Brand,
			CarModel:   mc.Criteria.CarModel,
			YearFrom:   mc.Criteria.KmFrom,
			YearTo:     mc.Criteria.YearTo,
			Fuel:       &mc.Criteria.Fuel,
			KmFrom:     mc.Criteria.KmFrom,
			KmTo:       mc.Criteria.KmTo,
			PageNumber: 1,
		}

		pageToScrapeJob := jobs.PageToScrapeJob{
			SessionID:  mc.SessionID,
			MarketID:   mc.MarketID,
			CriteriaID: mc.CriteriaID,
			PageURL:    pageToScrape,
			Visited:    false,
		}

		ptsBytes, err := json.Marshal(pageToScrapeJob)
		if err != nil {
			panic(err)
		}
		sms.messageQueue.PutMessage(sms.pagesToScrapeTopicName, ptsBytes)

	}
}

func (sms SessionManagerService) checkCompletedPages() {
	//
	message := sms.messageQueue.GetMessageWithDelete(sms.visitedPagesTopicName)
	var visitedPage jobs.PageToScrapeJob
	if len(*message) > 0 {
		err := json.Unmarshal(*message, &visitedPage)
		if err != nil {
			// push message back in the queue
			//sjc.messageQueue.PutMessage(sjc.pagesToScrapeTopicName, *message)
			panic(err)
		}

	}
}
