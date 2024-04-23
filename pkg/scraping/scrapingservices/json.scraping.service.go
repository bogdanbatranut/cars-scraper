package scrapingservices

import (
	"carscraper/pkg/jobs"
	"context"
	"log"
)

type JSONScrapingService struct {
	context        context.Context
	jobChannel     chan jobs.SessionJob
	resultsChannel chan jobs.AdsPageJobResult
	scrapingMapper IScrapingMapper
}

func NewJSONScrapingService(ctx context.Context, scrapingMapper IScrapingMapper) *JSONScrapingService {
	return &JSONScrapingService{
		context:        ctx,
		jobChannel:     make(chan jobs.SessionJob),
		resultsChannel: make(chan jobs.AdsPageJobResult, 1),
		scrapingMapper: scrapingMapper,
	}
}

func (service JSONScrapingService) Start() {
	log.Println("JSON Scraping Service Start")
	go func() {
		for {
			select {
			case job := <-service.jobChannel:
				go func() {
					service.processJob(job)
				}()
			case <-service.context.Done():
				log.Println("Colly Scraping Service Terminating")
				return
			}
		}
	}()
}

func (service JSONScrapingService) AddJob(job jobs.SessionJob) {
	service.jobChannel <- job
}

func (service JSONScrapingService) GetResultsChannel() *chan jobs.AdsPageJobResult {
	tmp := service.resultsChannel
	return &tmp
}

func (service JSONScrapingService) processJob(job jobs.SessionJob) {
	adapter := service.scrapingMapper.GetJSONMarketAdsAdapter(job.Market.Name)
	results := adapter.GetAds(job)
	adResult := jobs.AdsPageJobResult{
		RequestedScrapingJob: job,
		PageNumber:           job.Market.PageNumber,
		IsLastPage:           results.IsLastPage,
		Success:              results.Error == nil,
		Data:                 results.Ads,
	}
	if results.Ads == nil {
		return
	}
	log.Println("Colly service DONE job . Found ", len(*results.Ads), "ads")
	go func() {
		service.pushResultsToChannel(adResult)
	}()

	log.Println("JSON pushed results to channel")
}

func (service JSONScrapingService) pushResultsToChannel(res jobs.AdsPageJobResult) {
	service.resultsChannel <- res
}
