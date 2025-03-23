package tiriacauto

import (
	"carscraper/pkg/jobs"
	"fmt"
)

type TiriacAutoURLBuilder struct {
	brandsMap map[string]string
	modelsMap map[string]string
	fuelsMap  map[string]string
}

func initBrandNamesMap() map[string]string {
	brandsMap := make(map[string]string)
	brandsMap["bmw"] = "bmw"
	brandsMap["mercedes-benz"] = "mercedes-benz"
	brandsMap["volvo"] = "volvo"
	brandsMap["skoda"] = "skoda"
	brandsMap["opel"] = "opel"
	brandsMap["toyota"] = "toyota"
	brandsMap["volkswagen"] = "volkswagen"
	brandsMap["audi"] = "audi"
	brandsMap["mazda"] = "mazda"

	return brandsMap
}

func initModelsMap() map[string]string {
	modelsMap := make(map[string]string)

	modelsMap["x3"] = "x3"
	modelsMap["x3-m"] = "x3"
	modelsMap["x4"] = "x4"
	modelsMap["x4-m"] = "x4-m"
	modelsMap["x5"] = "x5"
	modelsMap["x5-m"] = "x5-m"
	modelsMap["x6"] = "x6"
	modelsMap["x6-m"] = "x6-m"
	modelsMap["7-series"] = "seria-7"
	modelsMap["gle-class"] = "gle"
	modelsMap["e-class"] = "e"
	modelsMap["s90"] = "s90"
	modelsMap["xc90"] = "xc90"
	modelsMap["xc60"] = "xc60"
	modelsMap["xc40"] = "xc40"
	modelsMap["glb-class"] = "glb"
	modelsMap["x3"] = "x3"
	modelsMap["glc-class"] = "glc"
	modelsMap["glc-coupe"] = "glc-coupe"
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
	modelsMap["cx-60"] = "cx-60"
	modelsMap["kodiaq"] = "kodiaq"
	return modelsMap
}

func initFuelsMap() map[string]string {
	fuelMap := make(map[string]string)
	fuelMap["diesel"] = "diesel"
	fuelMap["petrol"] = "gas"
	fuelMap["hybrid-petrol"] = "Hybrid+Plug-In+(Benzina)"
	fuelMap["hybrid-diesel"] = "Hybrid+Plug-In+(Diesel)"
	fuelMap["hybrid"] = "hibrid"
	fuelMap["other"] = "other"
	return fuelMap
}

func NewTiriacAutoURLBuilder() *TiriacAutoURLBuilder {
	return &TiriacAutoURLBuilder{
		brandsMap: initBrandNamesMap(),
		modelsMap: initModelsMap(),
		fuelsMap:  initFuelsMap(),
	}
}

func (b TiriacAutoURLBuilder) GetURL(job jobs.SessionJob) *string {
	brand := b.brandsMap[job.Criteria.Brand]
	model := b.modelsMap[job.Criteria.CarModel]
	fuel := b.fuelsMap[job.Criteria.Fuel]
	maxKm := *job.Criteria.KmTo
	minYear := job.Criteria.YearFrom

	url := fmt.Sprintf("https://tiriacauto.ro/auto-noi--auto-rulate/%s/%s?an-de-la=%d&combustibil=%s&km-pana-la=%d", brand, model, *minYear, fuel, maxKm)
	return &url
}
