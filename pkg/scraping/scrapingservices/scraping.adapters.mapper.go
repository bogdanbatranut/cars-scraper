package scrapingservices

import (
	"carscraper/pkg/scraping/marketadapters"
	"carscraper/pkg/scraping/markets/autoscout"
	"carscraper/pkg/scraping/markets/autotrack"
	"carscraper/pkg/scraping/markets/autovit"
	"carscraper/pkg/scraping/markets/mobile"
)

type IScrapingMapper interface {
	GetCollyMarketAdsAdapter(marketName string) marketadapters.ICollyMarketAdsAdapter
	GetRodMarketAdsAdapter(marketName string) marketadapters.IRodMarketAdsAdapter
	GetJSONMarketAdsAdapter(marketName string) marketadapters.IJSONMarketAdsAdapter
}

type ScrapingAdaptersMapper struct {
	collyAdapters map[string]marketadapters.ICollyMarketAdsAdapter
	rodAdapters   map[string]marketadapters.IRodMarketAdsAdapter
	jsonAdapters  map[string]marketadapters.IJSONMarketAdsAdapter
}

func NewScrapingAdaptersMapper() *ScrapingAdaptersMapper {

	collyAdaptersMap := make(map[string]marketadapters.ICollyMarketAdsAdapter)
	rodAdaptersMap := make(map[string]marketadapters.IRodMarketAdsAdapter)
	jsonAdaptersMap := make(map[string]marketadapters.IJSONMarketAdsAdapter)

	collyAdaptersMap["mobile.de"] = mobile.NewMobileDECollyMarketAdapter()
	rodAdaptersMap["autotracknl"] = autotrack.NewAutoTrackNLRodAdapter()
	rodAdaptersMap["autoscout"] = autoscout.NewAutoscoutRodAdapter()
	jsonAdaptersMap["autovit"] = autovit.NewAutovitJSONAdapter()

	return &ScrapingAdaptersMapper{
		collyAdapters: collyAdaptersMap,
		rodAdapters:   rodAdaptersMap,
		jsonAdapters:  jsonAdaptersMap,
	}
}

func (sm ScrapingAdaptersMapper) GetCollyMarketAdsAdapter(marketName string) marketadapters.ICollyMarketAdsAdapter {
	adapter := sm.collyAdapters[marketName]
	return adapter
}

func (sm ScrapingAdaptersMapper) GetRodMarketAdsAdapter(marketName string) marketadapters.IRodMarketAdsAdapter {
	adapter := sm.rodAdapters[marketName]
	return adapter
}

func (sm ScrapingAdaptersMapper) GetJSONMarketAdsAdapter(marketName string) marketadapters.IJSONMarketAdsAdapter {
	adapter := sm.jsonAdapters[marketName]
	return adapter
}
