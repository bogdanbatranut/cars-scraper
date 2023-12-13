package markets

import (
	"carscraper/pkg/scraping/markets/autovit"
	"carscraper/pkg/scraping/markets/mobile"
)

type ImplementationStrategies struct {
	strategies map[string]IScrapingStrategy
}

func NewImplemetationStrategies() ImplementationStrategies {
	s := make(map[string]IScrapingStrategy)
	// here we add to the map the implementations ...
	s["autovit"] = autovit.NewAutovitStrategy()
	s["mobile.de"] = mobile.NewMobileDeStrategy()
	//s["webcar"] = webcar.NewWebCarStrategy()
	is := ImplementationStrategies{
		strategies: s,
	}
	return is
}

func (is ImplementationStrategies) GetImplementation(marketName string) IScrapingStrategy {
	return is.strategies[marketName]
}
