package scrapingservices

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/logging"
	"carscraper/pkg/scraping/marketadapters"
	"carscraper/pkg/scraping/markets/autoklass"
	"carscraper/pkg/scraping/markets/autoscout"
	"carscraper/pkg/scraping/markets/autotrack"
	"carscraper/pkg/scraping/markets/autovit"
	"carscraper/pkg/scraping/markets/bmw_de"
	mercedes_benz_de "carscraper/pkg/scraping/markets/mercedes-benz_de"
	mercedes_benz_ro "carscraper/pkg/scraping/markets/mercedes-benz_ro"
	"carscraper/pkg/scraping/markets/mobile"
	"carscraper/pkg/scraping/markets/oferte_bmw"
	"carscraper/pkg/scraping/markets/olx"
	"carscraper/pkg/scraping/markets/tiriacauto"
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
	cfg           amconfig.IConfig
	//loggingService *logging.ScrapeLoggingService
}

func NewScrapingAdaptersMapper(loggingService *logging.ScrapeLoggingService) *ScrapingAdaptersMapper {

	//cfg, err := amconfig.NewViperConfig()
	//errorshandler.HandleErr(err)

	collyAdaptersMap := make(map[string]marketadapters.ICollyMarketAdsAdapter)
	rodAdaptersMap := make(map[string]marketadapters.IRodMarketAdsAdapter)
	jsonAdaptersMap := make(map[string]marketadapters.IJSONMarketAdsAdapter)

	collyAdaptersMap["mobile.de"] = mobile.NewMobileDECollyMarketAdapter(loggingService)
	//collyAdaptersMap["mobile.de"] = mobile_props.NewMobilePropDECollyMarketAdapter(loggingService, cfg)
	//collyAdaptersMap["mobile.de"] = mobilede_de.NewMobileDE_DECollyMarketAdapter(loggingService, cfg)
	collyAdaptersMap["bmw.de"] = bmw_de.NewBMWDECollyMarketAdapter(loggingService)
	collyAdaptersMap["tiriacauto"] = tiriacauto.NewTiriacAutoCollyMarketAdapter(loggingService)

	rodAdaptersMap["autotracknl"] = autotrack.NewAutoTrackNLRodAdapter()
	rodAdaptersMap["autoscout"] = autoscout.NewAutoscoutRodAdapter()

	jsonAdaptersMap["autovit"] = autovit.NewAutovitJSONAdapter(loggingService)
	jsonAdaptersMap["olx"] = olx.NewOLXJSONAdapter(loggingService)
	jsonAdaptersMap["oferte_bmw"] = oferte_bmw.NewOferteBMWAdapter(loggingService)
	jsonAdaptersMap["mercedes-benz.ro"] = mercedes_benz_ro.NewMercedesBenzRoAdapter(loggingService)
	jsonAdaptersMap["mercedes-benz.de"] = mercedes_benz_de.NewMercedesBenzDEAdapter(loggingService)
	jsonAdaptersMap["autoklass.ro"] = autoklass.NewAutoklassRoAdapter(loggingService)

	return &ScrapingAdaptersMapper{
		collyAdapters: collyAdaptersMap,
		rodAdapters:   rodAdaptersMap,
		jsonAdapters:  jsonAdaptersMap,
	}
}

func (sm ScrapingAdaptersMapper) SetCFG(cfg amconfig.IConfig) {
	sm.cfg = cfg
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
