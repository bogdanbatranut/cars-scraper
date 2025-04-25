package scrapingservices

import (
	"carscraper/pkg/scraping/icollector"
)

type CollyCollectorMapper struct {
	collectorsMap map[string]icollector.ICollyCollector
}

func NewCollyCollectorMapper() *CollyCollectorMapper {
	return &CollyCollectorMapper{
		collectorsMap: make(map[string]icollector.ICollyCollector),
	}
}

func (mapper CollyCollectorMapper) AddMarketCollyCollector(market string, marketCollyCollector icollector.ICollyCollector) {
	mapper.collectorsMap[market] = marketCollyCollector
}

func (mapper CollyCollectorMapper) GetCollector(market string) icollector.ICollyCollector {
	return mapper.GetCollector(market)
}
