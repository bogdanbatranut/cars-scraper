package url

import (
	"carscraper/pkg/adsdb"
	"fmt"
)

type AutovitURLComposer struct {
}

func NewAutovitURLComposer() *AutovitURLComposer {
	return &AutovitURLComposer{}
}

func (ac AutovitURLComposer) Create(criteria adsdb.Criteria) string {
	pageURL := fmt.Sprintf("https://www.autovit.ro/autoturisme/%s/%s/de-la-%d?search%%5Bfilter_enum_fuel_type%%5D=%s&search%%5Bfilter_float_mileage%%3Ato%%5D=%d",
		criteria.Brand,
		criteria.Model,
		criteria.YearFrom,
		criteria.Fuel,
		criteria.KmTo)
	//if sc.PageNumber != nil {
	//	pageURL = fmt.Sprintf("%s&page=%d", pageURL, *sc.PageNumber)
	//}
	return pageURL
}
