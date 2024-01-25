package markets

import (
	"carscraper/pkg/logging"
	"carscraper/pkg/scraping/markets/autoscout"
	"carscraper/pkg/scraping/markets/autotrack"
	"carscraper/pkg/scraping/markets/autovit"
	"carscraper/pkg/scraping/markets/mobile"
	"carscraper/pkg/scraping/markets/olx"
)

type ImplementationStrategies struct {
	strategies map[string]IScrapingStrategy
}

func NewImplemetationStrategies(logger logging.ScrapeLoggingService) ImplementationStrategies {
	s := make(map[string]IScrapingStrategy)
	// here we add to the map the implementations ...
	//s["webcar"] = webcar.NewWebCarStrategy()

	s["autovit"] = autovit.NewAutovitStrategy(logger)
	s["mobile.de"] = mobile.NewMobileDeStrategy(logger)
	s["autoscout"] = autoscout.NewAutoscoutStrategy(logger)
	s["autotracknl"] = autotrack.NewAutoTrackStrategy(&logger)
	s["olx"] = olx.NewOlxStrategy(&logger)
	is := ImplementationStrategies{
		strategies: s,
	}
	return is
}

func (is ImplementationStrategies) GetImplementation(marketName string) IScrapingStrategy {
	return is.strategies[marketName]
}
