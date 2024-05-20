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
	modelsMap map[string]string
	fuelsMap  map[string]string
}

func NewURLBuilder() *URLBuilder {
	return &URLBuilder{
		modelsMap: initModelsAdapterMap(),
		fuelsMap:  initFuelsMap(),
	}
}

func (b URLBuilder) GetURL(job jobs.SessionJob) *string {
	brand := job.Criteria.Brand
	model := b.modelsMap[job.Criteria.CarModel]
	fuel := b.fuelsMap[job.Criteria.Fuel]
	cyParam := url.QueryEscape("D,A,B,E,F,I,L,NL")
	ustateParam := url.QueryEscape("N,U")
	// https://www.autoscout24.ro/lst/mercedes-benz/gle-(toate)/ft_motorina?atype=C&cy=D%2CA%2CB%2CE%2CF%2CI%2CL%2CNL&desc=0&fregfrom=2019&kmto=125000&powertype=kw&search_id=21p91bbp3zv&sort=standard&source=detailsearch&ustate=N%2CU
	//https://www.autoscout24.ro/lst/mercedes-benz/glc-(toate)/ft_motorina?atype=C&cy=D%2CA%2CB%2CE%2CF%2CI%2CL%2CNL&damaged_listing=exclude&desc=0&fregfrom=2019&kmto=125000&powertype=kw&regfrom=2019&search_id=1kf3w3r2bjf&sort=price&source=detailsearch&ustate=N%2CU
	url := fmt.Sprintf("https://www.autoscout24.ro/lst/%s/%s/%s?atype=C&cy=%s&damaged_listing=exclude&desc=0&fregfrom=%d&kmto=%d&page=%d&powertype=kw&regfrom=%d&sort=price&source=detailsearch&ustate=%s", brand, model, fuel, cyParam, *job.Criteria.YearFrom, *job.Criteria.KmTo, job.Market.PageNumber, *job.Criteria.YearFrom, ustateParam)
	return &url
}

func initFuelsMap() map[string]string {
	fuelsMap := make(map[string]string)
	fuelsMap["diesel"] = "ft_motorina"
	fuelsMap["petrol"] = "ft_benzină"
	fuelsMap["hybrid-petrol"] = "ft_electric%2Fbenzină"
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
	modelsMap["xc40"] = "xc40"
	modelsMap["x3"] = "x3"
	modelsMap["glb-class"] = "glb-(toate)"
	modelsMap["glc-class"] = "glc-(toate)"
	modelsMap["glc-coupe"] = "glc-(toate)" // TODO implement this special case
	modelsMap["octavia"] = "octavia"
	modelsMap["superb"] = "superb"
	modelsMap["mokka"] = "mokka"
	modelsMap["yaris-cross"] = "yaris-cross"
	modelsMap["touareg"] = "touareg"
	modelsMap["a6"] = "a6"
	modelsMap["q8"] = "q8"
	modelsMap["q7"] = "q7"
	modelsMap["q5"] = "q5"
	modelsMap["q3"] = "q3"

	return modelsMap
}

// https://www.autoscout24.ro/lst/toyota//ft_benzină?atype=C&cy=D%2CA%2CB%2CE%2CF%2CI%2CL%2CNL&damaged_listing=exclude&desc=0&fregfrom=2019&kmto=125000&page=3&powertype=kw&regfrom=2019&sort=price&source=detailsearch&ustate=N%2CU
// https://www.autoscout24.ro/lst/toyota/yaris-cross/ft_benzin%C4%83?atype=C&cy=D%2CA%2CB%2CE%2CF%2CI%2CL%2CNL&damaged_listing=exclude&desc=0&fregfrom=2019&kmto=125000&powertype=kw&regfrom=2019&search_id=10tl4d5e8eu&sort=price&source=detailsearch&ustate=N%2CU&zipr=200
