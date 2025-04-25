package mercedes_benz_de

import (
	"carscraper/pkg/jobs"
)

type RequestBuilder struct {
	namingMapper *MercedesBenzDENamingMapper
}

func NewRequestBuilder() *RequestBuilder {
	namingMapper := NewMercedesBenzDENamingMapper()
	return &RequestBuilder{namingMapper: namingMapper}
}

func (b RequestBuilder) GetRequestBody(job jobs.SessionJob) Search {
	pageIndex := job.Market.PageNumber - 1
	return b.newBuildRequest(job.Criteria, pageIndex)
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
			Locale:    "de_DE",
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
