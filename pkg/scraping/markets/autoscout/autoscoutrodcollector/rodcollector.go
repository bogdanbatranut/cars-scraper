package autoscoutrodcollector

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
	"carscraper/pkg/scraping/scrapingservices"
)

type AutoScoutRodCollector struct {
	ResultsChannel   chan icollector.AdsResults
	RodPageProcessor *icollector.IRodPageProcessor
	RodService       *scrapingservices.RodBrowserService
}

func NewAutoScoutRodCollector(rodService *scrapingservices.RodBrowserService, processor *icollector.IRodPageProcessor) *AutoScoutRodCollector {
	return &AutoScoutRodCollector{
		ResultsChannel:   make(chan icollector.AdsResults),
		RodPageProcessor: processor,
		RodService:       rodService,
	}
}

func (collector AutoScoutRodCollector) GetAds(url string, pageNumber int, criteria jobs.Criteria) icollector.AdsResults {

	//collector.RodBrowserService.AddJob(url, *collector.RodPageProcessor)

	return icollector.AdsResults{
		Ads:        nil,
		IsLastPage: true,
		Error:      nil,
	}
}
