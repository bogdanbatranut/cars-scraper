package scrapingservices

type IScrapingServicesMapper interface {
	GetScrapingService(marketName string) IScrapingService
	AddMarketService(market string, service IScrapingService)
	GetAllServices() []IScrapingService
}

type ScrapingServicesMapper struct {
	scrapingServices map[string]IScrapingService
}

func NewScrapingServicesMapper() *ScrapingServicesMapper {
	mapper := ScrapingServicesMapper{
		scrapingServices: make(map[string]IScrapingService),
	}
	return &mapper
}

func (mapper ScrapingServicesMapper) AddMarketService(market string, service IScrapingService) {
	mapper.scrapingServices[market] = service
}

func (mapper ScrapingServicesMapper) GetScrapingService(market string) IScrapingService {
	return mapper.scrapingServices[market]
}

func (mapper ScrapingServicesMapper) GetAllServices() []IScrapingService {
	var services []IScrapingService
	for _, service := range mapper.scrapingServices {
		services = append(services, service)
	}
	return services
}
