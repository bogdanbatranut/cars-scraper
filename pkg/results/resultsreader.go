package results

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/events"
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/notifications"
	"carscraper/pkg/repos"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ResultsConsumerService struct {
	messageQueue        repos.IMessageQueue
	resultsTopicName    string
	scrapeResults       SessionCriteriaMarketResultsHandler
	pageAdsChannel      chan jobs.AdsPageJobResult
	resultsWriter       ResultsWriter
	logger              *logging.ScrapeLoggingService
	repo                *repos.AdsRepository
	notificationService *notifications.NotificationsService
	eventsListener      *events.EventsListener
	resultsManager      *ResultsManager
}

type ResultsReaderServiceConfiguration func(rcs *ResultsConsumerService)

func NewResultsReaderService(cfgs ...ResultsReaderServiceConfiguration) *ResultsConsumerService {
	res := NewSessionCriteriaMarketResults()
	service := &ResultsConsumerService{
		scrapeResults:  *res,
		pageAdsChannel: make(chan jobs.AdsPageJobResult),
		resultsManager: NewResultsManager(),
	}
	for _, cfg := range cfgs {
		cfg(service)
	}
	return service
}

func WithEventsListener(eventsListener *events.EventsListener) ResultsReaderServiceConfiguration {
	return func(cis *ResultsConsumerService) {
		cis.eventsListener = eventsListener
	}
}

func WithNotificationService(notificationService *notifications.NotificationsService) ResultsReaderServiceConfiguration {
	return func(cis *ResultsConsumerService) {
		cis.notificationService = notificationService
	}
}

func WithLogger(cfg amconfig.IConfig) ResultsReaderServiceConfiguration {
	logger := logging.NewScrapeLoggingService(cfg)
	return func(cis *ResultsConsumerService) {
		cis.logger = logger
	}
}

func WithResultsMQRepository(cfg amconfig.IConfig) ResultsReaderServiceConfiguration {
	smqHost := cfg.GetString(amconfig.SMQURL)
	smqPort := cfg.GetString(amconfig.SMQHTTPPort)
	mqHost := fmt.Sprintf("http://%s:%s", smqHost, smqPort)
	log.Println("MQ HOST: ", mqHost)
	smqr := repos.NewSimpleMessageQueueRepository(mqHost)
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

func WithRepo(repo *repos.AdsRepository) ResultsReaderServiceConfiguration {
	return func(rcs *ResultsConsumerService) {
		rcs.repo = repo
	}
}

func (rcs ResultsConsumerService) Start() {
	log.Println("Results Consumer Service StartAsync")
	done := make(chan bool, 1)

	signalsChannel := make(chan os.Signal, 1)
	signal.Notify(signalsChannel, syscall.SIGINT, syscall.SIGTERM)
	log.Println("start waiting for signal")

	_, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			time.Sleep(3 * time.Second)
			rcs.getResultsFromMQ()
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

		criteriaLog, err := rcs.logger.GetCriteriaLog(result.RequestedScrapingJob.SessionID, result.RequestedScrapingJob.CriteriaID, result.RequestedScrapingJob.MarketID)
		if err != nil {
			log.Println(err)
		}

		pageLog := rcs.logger.GetPageLog(result.RequestedScrapingJob.SessionID,
			result.RequestedScrapingJob.JobID,
			criteriaLog.ID,
			result.RequestedScrapingJob.MarketID,
			result.RequestedScrapingJob.Market.PageNumber,
		)

		if !result.Success {
			// TODO implement error in result...
			rcs.logger.PageLogSetError(pageLog, "NOT SUCCESSFUL")
			continue
		}

		if result.Data != nil {
			log.Printf("Got %d ads for market: %s ==> %s %s", len(*result.Data), result.RequestedScrapingJob.Market.Name, result.RequestedScrapingJob.Criteria.Brand, result.RequestedScrapingJob.Criteria.CarModel)
		} else {
			log.Println("The found ads are nil")
		}

		// TODO results.Data might be null... this happens on mobile when traversing pages... at some point you just get an empty page..

		if result.Data == nil || len(*result.Data) == 0 {
			//log.Printf("results %+v", result)
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
		// add results to the scrapeResults
		sessionIDStr := result.RequestedScrapingJob.SessionID.String()
		rcs.resultsManager.AddPageResults(sessionIDStr, result.RequestedScrapingJob.CriteriaID, result.RequestedScrapingJob.MarketID, result)
		err = rcs.logger.PageLogSetConsumed(pageLog)
		if err != nil {
			log.Println(err.Error())
		}
		rcs.logger.CriteriaLogAddNumberOfAds(*criteriaLog, len(*result.Data))
		completeCriteriaInMarket := rcs.resultsManager.isCompleteMarket(sessionIDStr, result.RequestedScrapingJob.CriteriaID, result.RequestedScrapingJob.MarketID)
		isCompleteCriteria := rcs.resultsManager.isCompleteCriteria(sessionIDStr, result.RequestedScrapingJob.CriteriaID)
		if completeCriteriaInMarket {
			brand := result.RequestedScrapingJob.Criteria.Brand
			model := result.RequestedScrapingJob.Criteria.CarModel
			market := result.RequestedScrapingJob.Market.Name
			marketID := result.RequestedScrapingJob.MarketID
			criteriaID := result.RequestedScrapingJob.CriteriaID

			ads := rcs.resultsManager.GetAds(sessionIDStr, result.RequestedScrapingJob.CriteriaID, result.RequestedScrapingJob.MarketID)
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
					//log.Printf("Deleted record with ID: %d", exsitingAdID)
				}
			}
			rcs.logger.CriteriaLogSetAsFinished(*criteriaLog)
			rcs.logger.CriteriaLogSetSuccessful(*criteriaLog)

			//TODO min price updated needs more complex logic
			//// get all ads in criteria
			//
			//var adsInCriteria []adsdb.Ad
			//tx := rcs.repo.GetDB().Model(&adsdb.Ad{}).Preload("Prices").Find(&adsInCriteria, *exsitingAdsIDs)
			//if tx.Error != nil {
			//	log.Println("Error getting ads in criteria")
			//}
			//minDbPrice := 10000000
			//var minPriceAd adsdb.Ad
			//for _, ad := range adsInCriteria {
			//	lastPrice := ad.Prices[len(ad.Prices)-1].Price
			//	if minDbPrice < lastPrice {
			//		minPriceAd = ad
			//	}
			//}
			//rcs.eventsListener.Fire(events.MinPriceUpdatedEvent{Ad: minPriceAd})
		}
		if isCompleteCriteria {
			//// get min price for ad in criteria
			//todayCheapestAd, err := rcs.repo.GetTodaysLowestPriceInCriteria(result.RequestedScrapingJob.CriteriaID)
			//if err != nil {
			//	log.Println("Error getting today's lowest price in criteria")
			//}
			//if todayCheapestAd != nil {
			//	err = rcs.notificationService.SendNewMinPrice(*todayCheapestAd)
			//	if err != nil {
			//		log.Println("Error sending new min) price notification")
			//	}
			//}
			//
			//todaysNewCheapestAd, err := rcs.repo.GetNewEntryLowestPrice(result.RequestedScrapingJob.CriteriaID)
			//if err != nil {
			//	log.Println("Error getting today's lowest price in criteria")
			//}
			//if todaysNewCheapestAd != nil {
			//	err = rcs.notificationService.SendNewMinPrice(*todaysNewCheapestAd)
			//	if err != nil {
			//		log.Println("Error sending new min price notification")
			//	}
			//
			//}

			//cheapestInCriteria, err := rcs.repo.GetMinPriceForCriteria(result.RequestedScrapingJob.CriteriaID)
			//if err != nil {
			//	log.Println("Error getting today's lowest price in criteria")
			//}
			//for _, ad := range cheapestInCriteria {
			//	err = rcs.notificationService.SendMinPriceInCriteria(ad)
			//	if err != nil {
			//		log.Println("Error sending new min price notification")
			//	}
			//}

		}
	}
}

//func (rcs ResultsConsumerService) processResultsOld() {
//	for {
//		result := <-rcs.pageAdsChannel
//
//		criteriaLog, err := rcs.logger.GetCriteriaLog(result.RequestedScrapingJob.SessionID, result.RequestedScrapingJob.CriteriaID, result.RequestedScrapingJob.MarketID)
//		if err != nil {
//			log.Println(err)
//		}
//
//		pageLog := rcs.logger.GetPageLog(result.RequestedScrapingJob.SessionID,
//			result.RequestedScrapingJob.JobID,
//			criteriaLog.ID,
//			result.RequestedScrapingJob.MarketID,
//			result.RequestedScrapingJob.Market.PageNumber,
//		)
//
//		if !result.Success {
//			// TODO implement error in result...
//			rcs.logger.PageLogSetError(pageLog, "NOT SUCCESSFUL")
//			continue
//		}
//
//		if result.Data != nil {
//			log.Printf("Got %d ads for market: %s ==> %s %s", len(*result.Data), result.RequestedScrapingJob.Market.Name, result.RequestedScrapingJob.Criteria.Brand, result.RequestedScrapingJob.Criteria.CarModel)
//
//		} else {
//			log.Println("The found ads are nil")
//		}
//
//		// TODO results.Data might be null... this happens on mobile when traversing pages... at some point you just get an empty page..
//
//		if result.Data == nil || len(*result.Data) == 0 {
//			//log.Printf("results %+v", result)
//			result.Data = &[]jobs.Ad{jobs.Ad{
//				Brand:              "",
//				Model:              "",
//				Year:               0,
//				Km:                 0,
//				Fuel:               "",
//				Price:              0,
//				AdID:               "",
//				Ad_url:             "",
//				SellerType:         "",
//				SellerName:         nil,
//				SellerNameInMarket: nil,
//				SellerOwnURL:       nil,
//				SellerMarketURL:    nil,
//			}}
//			result.IsLastPage = true
//			//continue
//		}
//
//		rcs.scrapeResults.Add(result.RequestedScrapingJob.SessionID, result.RequestedScrapingJob.CriteriaID, result.RequestedScrapingJob.MarketID, result)
//		err = rcs.logger.PageLogSetConsumed(pageLog)
//		if err != nil {
//			log.Println(err.Error())
//		}
//		rcs.logger.CriteriaLogAddNumberOfAds(*criteriaLog, len(*result.Data))
//		completeCriteriaInMarket := rcs.scrapeResults.results[result.RequestedScrapingJob.SessionID.String()][result.RequestedScrapingJob.CriteriaID][result.RequestedScrapingJob.MarketID].IsComplete()
//		if completeCriteriaInMarket {
//			brand := result.RequestedScrapingJob.Criteria.Brand
//			model := result.RequestedScrapingJob.Criteria.CarModel
//			market := result.RequestedScrapingJob.Market.Name
//			marketID := result.RequestedScrapingJob.MarketID
//			criteriaID := result.RequestedScrapingJob.CriteriaID
//
//			marketSrapintResults := rcs.scrapeResults.results[result.RequestedScrapingJob.SessionID.String()][result.RequestedScrapingJob.CriteriaID][result.RequestedScrapingJob.MarketID]
//			ads := marketSrapintResults.getAds()
//			totalAds := len(ads)
//			if ads == nil {
//				continue
//			}
//			log.Printf("WE HAVE A COMPLETE CRITERIA IN THE MARKET -> Brand: %s Model: %s Market: %s Total Ads: %d", brand, model, market, totalAds)
//			// transform them to db writeable results
//			upsertedAdsIDs, err := rcs.resultsWriter.WriteAds(ads, result.RequestedScrapingJob.MarketID, result.RequestedScrapingJob.CriteriaID)
//			if err != nil {
//				panic(err)
//			}
//			exsitingAdsIDs := rcs.resultsWriter.GetAllAdsIDs(marketID, criteriaID)
//
//			for _, exsitingAdID := range *exsitingAdsIDs {
//				found := false
//				for _, upsertedAdID := range *upsertedAdsIDs {
//					if exsitingAdID == upsertedAdID {
//						found = true
//						break
//					}
//				}
//				if !found {
//					rcs.resultsWriter.DeleteAd(exsitingAdID)
//					//log.Printf("Deleted record with ID: %d", exsitingAdID)
//				}
//			}
//			rcs.logger.CriteriaLogSetAsFinished(*criteriaLog)
//			rcs.logger.CriteriaLogSetSuccessful(*criteriaLog)
//
//			//TODO min price updated needs more complex logic
//			//// get all ads in criteria
//			//
//			//var adsInCriteria []adsdb.Ad
//			//tx := rcs.repo.GetDB().Model(&adsdb.Ad{}).Preload("Prices").Find(&adsInCriteria, *exsitingAdsIDs)
//			//if tx.Error != nil {
//			//	log.Println("Error getting ads in criteria")
//			//}
//			//minDbPrice := 10000000
//			//var minPriceAd adsdb.Ad
//			//for _, ad := range adsInCriteria {
//			//	lastPrice := ad.Prices[len(ad.Prices)-1].Price
//			//	if minDbPrice < lastPrice {
//			//		minPriceAd = ad
//			//	}
//			//}
//			//rcs.eventsListener.Fire(events.MinPriceUpdatedEvent{Ad: minPriceAd})
//		}
//	}
//}
