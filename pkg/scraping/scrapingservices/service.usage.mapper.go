package scrapingservices

type ServiceUsageMapper struct {
	servicesMap map[string]IScrapingService
}

func NewServiceUsageMapper() *ServiceUsageMapper {
	return &ServiceUsageMapper{}
}

func (sum ServiceUsageMapper) RegisterrScrapingService(market string, service IScrapingService) {
	sum.servicesMap[market] = service
}

func (sum ServiceUsageMapper) GetServiceForMarket(market string) IScrapingService {
	return sum.servicesMap[market]
}
