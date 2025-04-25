package markets

import (
	"carscraper/pkg/logging"
	"carscraper/pkg/scraping/markets/autoscout"
	"carscraper/pkg/scraping/markets/autoscout/autoscoutcollycollector"
	"carscraper/pkg/scraping/markets/autotrack"
	"carscraper/pkg/scraping/markets/autovit"
	"carscraper/pkg/scraping/markets/mobile"
	"carscraper/pkg/scraping/markets/olx"
	"carscraper/pkg/scraping/scrapingservices"
)

type ImplementationStrategies struct {
	strategies map[string]IScrapingJob
	//rodService *scrapingservices.RodBrowserService
}

func NewImplemetationStrategies(logger logging.ScrapeLoggingService, rodService *scrapingservices.RodBrowserService) ImplementationStrategies {
	s := make(map[string]IScrapingJob)
	// here we add to the map the implementations ...
	//s["webcar"] = webcar.NewWebCarStrategy()

	s["autovit"] = autovit.NewAutovitStrategy(logger)
	s["mobile.de"] = mobile.NewMobileDeStrategy(logger)

	autoscoutCollyCollector := autoscoutcollycollector.NewAutoscoutCollyProcessor()
	//s["autoscout"] = autoscout.NewAutoscoutStrategy(&logger, autoscoutCollyCollector, rodService)
	s["autoscout"] = autoscout.NewAutoscoutStrategy(&logger, autoscoutCollyCollector)

	s["autotracknl"] = autotrack.NewAutoTrackStrategy(&logger)
	s["olx"] = olx.NewOlxStrategy(&logger)
	is := ImplementationStrategies{
		strategies: s,
		//rodService: rodService,
	}
	return is
}

func (is ImplementationStrategies) GetImplementation(marketName string) IScrapingJob {
	return is.strategies[marketName]
}
