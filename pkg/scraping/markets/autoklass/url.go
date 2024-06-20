package autoklass

import (
	"carscraper/pkg/jobs"
	"fmt"
)

type AutoklassRoURLBuilder struct {
}

func NewAutoklassRoURLBuilder() *AutoklassRoURLBuilder {
	return &AutoklassRoURLBuilder{}
}

func (b AutoklassRoURLBuilder) GetURL(job jobs.SessionJob, namingMapper AutoklassRONamingMapper) string {
	year := *job.Criteria.YearFrom
	carBrandID := namingMapper.carBrandsIDs[job.Criteria.Brand]
	carModelID := namingMapper.carModelsIDs[job.Criteria.CarModel]
	fuel := namingMapper.fuelTypes[job.Criteria.Fuel]
	encodedSpace := "%20"
	url := fmt.Sprintf("https://www.autoklass.ro/api/cars/carSearch?budgetFrom=0&budgetTo=1000000&order=salePrice%sASC&yearManufactureTo=2038&yearManufactureFrom=%d&kmFrom=0&kmTo=250000&status=rulata&carBrandID=%d&carModelID=%d&limit=32&fuel=%s", encodedSpace, year, carBrandID, carModelID, fuel)
	return url
}
