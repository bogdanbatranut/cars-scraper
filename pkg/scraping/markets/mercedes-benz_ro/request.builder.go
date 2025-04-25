package mercedes_benz_ro

import (
	"carscraper/pkg/jobs"
)

type RequestBuilder struct {
	namingMapper *MercedesBenzRONamingMapper
}

func NewRequestBuilder() *RequestBuilder {
	namingMapper := NewMercedesBenzRoNamingMapper()
	return &RequestBuilder{namingMapper: namingMapper}
}

func (b RequestBuilder) GetRequestBody(job jobs.SessionJob) Search {
	pageIndex := job.Market.PageNumber - 1
	//yearFrom := *job.Criteria.YearFrom
	//yearFrom = yearFrom * 10000
	//yearFrom = yearFrom + 101
	//year, month, day := time.Now().Date()
	//yearToStr := fmt.Sprintf("%d%02d%02d", year, month, day)
	//yearTo, err := strconv.Atoi(yearToStr)
	//if err != nil {
	//	panic(err)
	//}
	return b.newBuildRequest(job.Criteria, pageIndex)
	//return b.buildRequest(job.Criteria.CarModel, job.Criteria.Fuel, yearFrom, yearTo, pageIndex)
}

func (b RequestBuilder) newBuildRequest(criteria jobs.Criteria, pageIndex int) Search {
	marketCriteria := b.namingMapper.GetMarketCriteria(criteria)

	vsr := VehicleSearchRequest{
		SearchInfo: SearchInfo{
			Paging: Paging{
				PageIndex: pageIndex,
				Quantity:  12,
			},
			Searchterm: Searchterm{FindCompleteTermOnly: true},
			Sort: []Sort{{
				Field: "offerPriceGross",
				Order: "ASC",
			}},
		},
		Facets:   []string{"salesClass"},
		Criteria: marketCriteria,
		Context: Context{
			ProcessId: "UCui",
			Locale:    "ro_RO",
			OutletIds: []interface{}{},
			UiId:      "main",
		},
	}
	req := Search{VehicleSearchRequest: vsr}

	return req
}

func (b RequestBuilder) buildRequest(model string, fuel string, yearFrom int, yearTo int, pageIndex int) Search {
	salesClass := b.namingMapper.GetSalesClassCodesTextEntry(model)
	fr := FirstRegistration{
		Max: yearTo,
		Min: yearFrom,
	}

	engineType := b.namingMapper.GetEngineTypeCodesTextEntry(fuel)

	vsr := VehicleSearchRequest{
		SearchInfo: SearchInfo{
			Paging: Paging{
				PageIndex: 0,
				Quantity:  1000,
			},
			Searchterm: Searchterm{FindCompleteTermOnly: true},
			Sort: []Sort{{
				Field: "offerPriceGross",
				Order: "ASC",
			}},
		},
		Facets: []string{"salesClass"},
		Criteria: Criteria{
			EngineType:        []CodesTextEntry{engineType},
			FirstRegistration: fr,
			SalesClass:        []CodesTextEntry{salesClass},
		},
		Context: Context{
			ProcessId: "UCui",
			Locale:    "ro_RO",
			OutletIds: []interface{}{},
			UiId:      "main",
		},
	}
	req := Search{VehicleSearchRequest: vsr}

	return req
}
