package mobilede

import (
	"carscraper/pkg/jobs"
	"fmt"
	"strings"
)

// https://www.mobile.de/
//ro/automobil/mercedes-benz-clasa-gle/vhc:car,pgn:1,pgs:50,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000,dmg:false
// url builder

type BrandModelValues struct {
	Brand    string
	Model    string
	SubModel *string
}

type QueryParam struct {
	Name  string
	Value string
}

func (q QueryParam) toQueryStr() string {
	return fmt.Sprintf("%s:%s")
}

type URLBuilder struct {
	criteria               jobs.Criteria
	modelsMap              map[string]string
	brandModelParamsValues map[string]map[string]BrandModelValues
	fuelParam              map[string]string
}

func NewURLBuilder(criteria jobs.Criteria) *URLBuilder {
	return &URLBuilder{
		criteria:               criteria,
		modelsMap:              initModelsAdapterMap(),
		brandModelParamsValues: initParamNames(),
		fuelParam:              initFuelParams(),
	}
}

func (b URLBuilder) GetMobileDEPageURL(pageNumber int) string {
	brand := b.criteria.Brand
	model := b.modelsMap[b.criteria.CarModel]
	brandParamValue := b.brandModelParamsValues[brand][model].Brand
	modelParamValue := b.brandModelParamsValues[brand][model].Model
	subModelParamValue := ""
	if b.brandModelParamsValues[brand][model].SubModel != nil {
		subModelParamValue = *b.brandModelParamsValues[brand][model].SubModel
	}

	fuelParam := ""
	fuel := b.fuelParam[b.criteria.Fuel]
	if fuel != "" {
		fuelParam = fmt.Sprintf(",ft:%s", strings.ToUpper(fuel))
	}

	brandModelParam := fmt.Sprintf("ms=%s;;%s;%s", brandParamValue, modelParamValue, subModelParamValue)

	url := fmt.Sprintf("https://suchen.mobile.de/fahrzeuge/search.html?dam=false&fr=%d:&%s&isSearchRequest=true&ml=:%d&%s&od=up&ps=0&ref=sroHead&s=Car&sb=p&vc=Car&pageNumber=%d&refId=801bb8bc-0d38-9218-3a38-95241f83e3ff", *b.criteria.YearFrom, fuelParam, *b.criteria.KmTo, brandModelParam, pageNumber)
	return url
}

func (b URLBuilder) GetPageURL(pageNumber int) string {
	brand := b.criteria.Brand
	model := b.modelsMap[b.criteria.CarModel]
	brandParamValue := b.brandModelParamsValues[brand][model].Brand
	modelParamValue := b.brandModelParamsValues[brand][model].Model
	subModelParamValue := ""
	if b.brandModelParamsValues[brand][model].SubModel != nil {
		subModelParamValue = *b.brandModelParamsValues[brand][model].SubModel
	}

	fuelParam := ""
	fuel := b.fuelParam[b.criteria.Fuel]
	if fuel != "" {
		fuelParam = fmt.Sprintf(",ful:%s", fuel)
	}

	url := fmt.Sprintf("https://www.mobile.de/ro/automobil/%s-%s/vhc:car,pgn:%d,pgs:50,srt:price,sro:asc,ms1:%s_%s_%s,frn:%d%s,mlx:125000,dmg:false", brand, model, pageNumber, brandParamValue, modelParamValue, subModelParamValue, *b.criteria.YearFrom, fuelParam)
	return url
}

func initFuelParams() map[string]string {
	fuelMap := make(map[string]string)
	fuelMap["diesel"] = "diesel"
	fuelMap["petrol"] = "petrol"
	fuelMap["hybrid-petrol"] = "hybrid"
	fuelMap["hybrid"] = "hybrid"
	return fuelMap
}

func initModelsAdapterMap() map[string]string {
	modelsMap := make(map[string]string)

	modelsMap["x4"] = "x4"
	modelsMap["x4-m"] = "x4-m"
	modelsMap["x5"] = "x5"
	modelsMap["x5-m"] = "x5-m"
	modelsMap["x6"] = "x6"
	modelsMap["x6-m"] = "x6-m"
	modelsMap["7-series"] = "seria-7"
	modelsMap["gle-class"] = "clasa-gle"
	modelsMap["e-class"] = "clasa-e"
	modelsMap["s90"] = "s90"
	modelsMap["xc90"] = "xc90"
	modelsMap["xc60"] = "xc60"
	modelsMap["xc40"] = "xc40"
	modelsMap["glb-class"] = "clasa-glb"
	modelsMap["x3"] = "x3"
	modelsMap["glc-class"] = "clasa-glc"
	modelsMap["glc-coupe"] = "clasa-glc-coupe"
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

func initParamNames() map[string]map[string]BrandModelValues {
	params := make(map[string]map[string]BrandModelValues)

	opModelsMap := map[string]BrandModelValues{}
	opMokka := BrandModelValues{
		Brand: "19000",
		Model: "37",
	}
	opModelsMap["mokka"] = opMokka

	skModelsMap := map[string]BrandModelValues{}
	sko := BrandModelValues{
		Brand: "22900",
		Model: "10",
	}
	skModelsMap["octavia"] = sko

	sks := BrandModelValues{
		Brand: "22900",
		Model: "12",
	}
	skModelsMap["superb"] = sks

	skodaKodiaq := BrandModelValues{
		Brand: "22900",
		Model: "19",
	}
	skModelsMap["kodiaq"] = skodaKodiaq

	mbModelsMap := map[string]BrandModelValues{}

	coupe := "coupe"
	mbGLC_Coupe := BrandModelValues{
		Brand:    "17200",
		Model:    "-59",
		SubModel: &coupe,
	}
	mbModelsMap["clasa-glc-coupe"] = mbGLC_Coupe

	mbGLC := BrandModelValues{
		Brand: "17200",
		Model: "-59",
	}
	mbModelsMap["clasa-glc"] = mbGLC

	mbGLB := BrandModelValues{
		Brand: "17200",
		Model: "-66",
	}
	mbModelsMap["clasa-glb"] = mbGLB

	mbGLE := BrandModelValues{
		Brand: "17200",
		Model: "-58",
	}
	mbModelsMap["clasa-gle"] = mbGLE

	mbE := BrandModelValues{
		Brand: "17200",
		Model: "-11",
	}
	mbModelsMap["clasa-e"] = mbE

	params["mercedes-benz"] = mbModelsMap
	params["opel"] = opModelsMap
	params["skoda"] = skModelsMap

	volvoModelsMap := map[string]BrandModelValues{}
	vs90 := BrandModelValues{
		Brand: "25100",
		Model: "31",
	}
	volvoModelsMap["s90"] = vs90

	vxc60 := BrandModelValues{
		Brand: "25100",
		Model: "40",
	}
	volvoModelsMap["xc60"] = vxc60

	vxc90 := BrandModelValues{
		Brand: "25100",
		Model: "37",
	}
	volvoModelsMap["xc90"] = vxc90

	vxc40 := BrandModelValues{
		Brand: "25100",
		Model: "45",
	}
	volvoModelsMap["xc40"] = vxc40

	bmwModelsMap := map[string]BrandModelValues{}
	params["volvo"] = volvoModelsMap

	bmwx3 := BrandModelValues{
		Brand: "3500",
		Model: "48",
	}
	bmwModelsMap["x3"] = bmwx3

	bmwx4 := BrandModelValues{
		Brand: "3500",
		Model: "92",
	}
	bmwModelsMap["x4"] = bmwx4

	bmwx4m := BrandModelValues{
		Brand: "3500",
		Model: "146",
	}
	bmwModelsMap["x4m"] = bmwx4m

	bmwx5 := BrandModelValues{
		Brand: "3500",
		Model: "49",
	}
	bmwModelsMap["x5"] = bmwx5

	bmwx5m := BrandModelValues{
		Brand: "3500",
		Model: "53",
	}
	bmwModelsMap["x5-m"] = bmwx5m

	bmwx6 := BrandModelValues{
		Brand: "3500",
		Model: "60",
	}
	bmwModelsMap["x6"] = bmwx6

	bmwx6m := BrandModelValues{
		Brand: "3500",
		Model: "62",
	}
	bmwModelsMap["x6-m"] = bmwx6m

	bmw7 := BrandModelValues{
		Brand: "3500",
		Model: "-24",
	}
	bmwModelsMap["seria-7"] = bmw7
	params["bmw"] = bmwModelsMap

	toyotaModelsMap := map[string]BrandModelValues{}
	toyotaYarisCross := BrandModelValues{
		Brand: "24100",
		Model: "78",
	}
	toyotaModelsMap["yaris-cross"] = toyotaYarisCross
	params["toyota"] = toyotaModelsMap

	vwModelsMap := map[string]BrandModelValues{}
	vwTouareg := BrandModelValues{
		Brand: "25200",
		Model: "36",
	}
	vwModelsMap["touareg"] = vwTouareg
	params["volkswagen"] = vwModelsMap

	audiModelsMap := map[string]BrandModelValues{}
	a6 := BrandModelValues{
		Brand: "1900",
		Model: "10",
	}
	audiModelsMap["a6"] = a6

	q8 := BrandModelValues{
		Brand: "1900",
		Model: "46",
	}
	audiModelsMap["q8"] = q8

	q7 := BrandModelValues{
		Brand: "1900",
		Model: "15",
	}
	audiModelsMap["q7"] = q7

	q5 := BrandModelValues{
		Brand: "1900",
		Model: "32",
	}
	audiModelsMap["q5"] = q5

	q3 := BrandModelValues{
		Brand: "1900",
		Model: "37",
	}
	audiModelsMap["q3"] = q3

	params["audi"] = audiModelsMap
	return params
}
