package scrapingservices

import (
	"carscraper/pkg/jobs"
	"context"
	"log"
)

type CollyScrapingService struct {
	context        context.Context
	jobChannel     chan jobs.SessionJob
	resultsChannel chan jobs.AdsPageJobResult
	scrapingMapper IScrapingMapper
}

func NewCollyScrapingService(ctx context.Context, scrapingMapper IScrapingMapper) *CollyScrapingService {
	return &CollyScrapingService{
		context:        ctx,
		jobChannel:     make(chan jobs.SessionJob),
		resultsChannel: make(chan jobs.AdsPageJobResult, 1),
		scrapingMapper: scrapingMapper,
	}
}

func (css CollyScrapingService) Start() {
	log.Println("Colly Scraping Service Start")
	go func() {
		for {
			select {
			case job := <-css.jobChannel:
				css.processJob(job)
			case <-css.context.Done():
				log.Println("Colly Scraping Service Terminating")
				return
			}
		}
	}()
}

func (css CollyScrapingService) AddJob(job jobs.SessionJob) {
	css.jobChannel <- job
}

func (css CollyScrapingService) GetResultsChannel() *chan jobs.AdsPageJobResult {
	tmp := css.resultsChannel
	return &tmp
}

func (css CollyScrapingService) processJob(job jobs.SessionJob) {
	adapter := css.scrapingMapper.GetCollyMarketAdsAdapter(job.Market.Name)
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
		css.pushResultsToChannel(adResult)
	}()

	log.Println("COLLY pushed results to channel")
}

func (css CollyScrapingService) pushResultsToChannel(res jobs.AdsPageJobResult) {
	css.resultsChannel <- res
}
