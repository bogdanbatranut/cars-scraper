package autoscout

import (
	"carscraper/pkg/jobs"
	"fmt"
	"net/url"
)

// https://www.mobile.de/
//ro/automobil/mercedes-benz-clasa-gle/vhc:car,pgn:1,pgs:50,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000,dmg:false
// url builder

type BrandModelValues struct {
	Brand string
	Model string
}

type QueryParam struct {
	Name  string
	Value string
}

func (q QueryParam) toQueryStr() string {
	return fmt.Sprintf("%s:%s")
}

type URLBuilder struct {
	criteria  jobs.Criteria
	modelsMap map[string]string
	fuelsMap  map[string]string
}

func NewURLBuilder(criteria jobs.Criteria) *URLBuilder {
	return &URLBuilder{
		criteria:  criteria,
		modelsMap: initModelsAdapterMap(),
		fuelsMap:  initFuelsMap(),
	}
}

func (b URLBuilder) GetPageURL(pageNumber int) string {
	brand := b.criteria.Brand
	model := b.modelsMap[b.criteria.CarModel]
	fuel := b.fuelsMap[b.criteria.Fuel]
	cyParam := url.QueryEscape("D,A,B,E,F,I,L,NL")
	ustateParam := url.QueryEscape("N,U")
	// https://www.autoscout24.ro/lst/mercedes-benz/gle-(toate)/ft_motorina?atype=C&cy=D%2CA%2CB%2CE%2CF%2CI%2CL%2CNL&desc=0&fregfrom=2019&kmto=125000&powertype=kw&search_id=21p91bbp3zv&sort=standard&source=detailsearch&ustate=N%2CU
	url := fmt.Sprintf("https://www.autoscout24.ro/lst/%s/%s/%s?atype=C&cy=%s&damaged_listing=exclude&desc=0&fregfrom=%d&kmto=%d&page=%d&powertype=kw&regfrom=%d&sort=standard&source=detailsearch&ustate=%s", brand, model, fuel, cyParam, *b.criteria.YearFrom, *b.criteria.KmTo, pageNumber, *b.criteria.YearFrom, ustateParam)
	return url
}

func initFuelsMap() map[string]string {
	fuelsMap := make(map[string]string)
	fuelsMap["diesel"] = "ft_motorina"
	return fuelsMap
}

func initModelsAdapterMap() map[string]string {
	modelsMap := make(map[string]string)

	modelsMap["x4"] = "x4"
	modelsMap["x4-m"] = "x4-m"
	modelsMap["x5"] = "x5"
	modelsMap["x5-m"] = "x5-m"
	modelsMap["x6"] = "x6"
	modelsMap["x6-m"] = "x6-m"
	modelsMap["7-series"] = "seria-7-(toate)"
	modelsMap["gle-class"] = "gle-(toate)" // TODO - this might be ok
	modelsMap["e-class"] = "clasa-e-(toate)"
	modelsMap["s90"] = "s90"
	modelsMap["xc90"] = "xc90"
	modelsMap["xc60"] = "xc60"

	return modelsMap
}
