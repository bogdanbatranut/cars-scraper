package results

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/jobs"
	"carscraper/pkg/repos"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type ResultsConsumerService struct {
	messageQueue     repos.IMessageQueue
	resultsTopicName string
	scrapeResults    SessionCriteriaMarketResultsHandler
	pageAdsChannel   chan jobs.AdsPageJobResult
	resultsWriter    ResultsWriter
}

type ResultsReaderServiceConfiguration func(rcs *ResultsConsumerService)

func NewResultsReaderService(cfgs ...ResultsReaderServiceConfiguration) *ResultsConsumerService {
	res := NewSessionCriteriaMarketResults()
	service := &ResultsConsumerService{
		scrapeResults:  *res,
		pageAdsChannel: make(chan jobs.AdsPageJobResult),
	}
	for _, cfg := range cfgs {
		cfg(service)
	}
	return service
}

func WithResultsMQRepository(cfg amconfig.IConfig) ResultsReaderServiceConfiguration {
	smqHost := cfg.GetString(amconfig.SMQURL)
	smqPort := cfg.GetString(amconfig.SMQHTTPPort)
	smqr := repos.NewSimpleMessageQueueRepository(fmt.Sprintf("http://%s:%s", smqHost, smqPort))
	return WithMessageQueueRepository(smqr)
}

func WithMessageQueueRepository(mqr repos.IMessageQueue) ResultsReaderServiceConfiguration {
	return func(cis *ResultsConsumerService) {
		cis.messageQueue = mqr
	}
}

func WithResultsWriter(rw ResultsWriter) ResultsReaderServiceConfiguration {
	return func(cis *ResultsConsumerService) {
		cis.resultsWriter = rw
	}
}

func WithTopicName(cfg amconfig.IConfig) ResultsReaderServiceConfiguration {
	topicName := cfg.GetString(amconfig.SMQResultsTopicName)
	return func(rcs *ResultsConsumerService) {
		rcs.resultsTopicName = topicName
	}
}

func (rcs ResultsConsumerService) Start() {
	log.Println("Results Consumer Service Start")
	done := make(chan bool, 1)

	signalsChannel := make(chan os.Signal, 1)
	signal.Notify(signalsChannel, syscall.SIGINT, syscall.SIGTERM)
	log.Println("start waiting for signal")

	_, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			rcs.getResultsFromMQ()
			//time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		rcs.processResults()
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

func (rcs ResultsConsumerService) getResultsFromMQ() {
	message := rcs.messageQueue.GetMessageWithDelete(rcs.resultsTopicName)
	var result jobs.AdsPageJobResult
	if len(*message) > 0 {
		err := json.Unmarshal(*message, &result)
		if err != nil {
			panic(err)
		}
		rcs.pageAdsChannel <- result
	}
}

func (rcs ResultsConsumerService) processResults() {
	for {
		result := <-rcs.pageAdsChannel
		if !result.Success {
			// TODO implement error in result...
			continue
		}

		// TODO results.Data might be null... this happens on mobile when traversing pages... at some point you just get an empty page..

		if result.Data == nil || len(*result.Data) == 0 {
			log.Printf("results %+v", result)
			result.Data = &[]jobs.Ad{jobs.Ad{
				Brand:              "",
				Model:              "",
				Year:               0,
				Km:                 0,
				Fuel:               "",
				Price:              0,
				AdID:               "",
				Ad_url:             "",
				SellerType:         "",
				SellerName:         nil,
				SellerNameInMarket: nil,
				SellerOwnURL:       nil,
				SellerMarketURL:    nil,
			}}
			result.IsLastPage = true
			//continue
		}

		rcs.scrapeResults.Add(result.RequestedScrapingJob.SessionID, result.RequestedScrapingJob.CriteriaID, result.RequestedScrapingJob.MarketID, result)
		//rcs.scrapeResults.Print()
		//log.Printf("Results for sessionID: %s Make: %s Model: %s TOTAL in MEM: %d",
		//	result.RequestedScrapingJob.SessionID.String(),
		//	result.RequestedScrapingJob.Criteria.Brand,
		//	result.RequestedScrapingJob.Criteria.CarModel,
		//	len(rcs.scrapeResults.results[result.RequestedScrapingJob.SessionID.String()][result.RequestedScrapingJob.CriteriaID][result.RequestedScrapingJob.MarketID].adsInPage))

		complete := rcs.scrapeResults.results[result.RequestedScrapingJob.SessionID.String()][result.RequestedScrapingJob.CriteriaID][result.RequestedScrapingJob.MarketID].IsComplete()
		if complete {
			brand := result.RequestedScrapingJob.Criteria.Brand
			model := result.RequestedScrapingJob.Criteria.CarModel
			market := result.RequestedScrapingJob.Market.Name
			marketID := result.RequestedScrapingJob.MarketID
			criteriaID := result.RequestedScrapingJob.CriteriaID

			marketSrapintResults := rcs.scrapeResults.results[result.RequestedScrapingJob.SessionID.String()][result.RequestedScrapingJob.CriteriaID][result.RequestedScrapingJob.MarketID]
			ads := marketSrapintResults.getAds()
			totalAds := len(ads)
			if ads == nil {
				continue
			}
			log.Printf("WE HAVE A COMPLETE CRITERIA IN THE MARKET -> Brand: %s Model: %s Market: %s Total Ads: %d", brand, model, market, totalAds)
			// transform them to db writeable results
			upsertedAdsIDs, err := rcs.resultsWriter.WriteAds(ads, result.RequestedScrapingJob.MarketID, result.RequestedScrapingJob.CriteriaID)
			if err != nil {
				panic(err)
			}
			exsitingAdsIDs := rcs.resultsWriter.GetAllAdsIDs(marketID, criteriaID)

			for _, exsitingAdID := range *exsitingAdsIDs {
				found := false
				for _, upsertedAdID := range *upsertedAdsIDs {
					if exsitingAdID == upsertedAdID {
						found = true
						break
					}
				}
				if !found {
					rcs.resultsWriter.DeleteAd(exsitingAdID)
					log.Printf("Deleted record with ID: %d", exsitingAdID)
				}
			}

		}
	}
}
