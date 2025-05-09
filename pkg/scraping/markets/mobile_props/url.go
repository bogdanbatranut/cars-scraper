package mobile_props

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/jobs"
	"carscraper/pkg/repos"
	"fmt"
)

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
	modelsMap              map[string]string
	brandModelParamsValues map[string]map[string]BrandModelValues
	fuelParam              map[string]string
	propertiesRepo         *repos.PropertiesRepository
	marketsRepo            *repos.SQLMarketsRepository
	marketID               uint
	propertiesRepository   *repos.PropertiesRepository
}

func NewURLBuilder(cfg amconfig.IConfig) *URLBuilder {
	marketsRepo := repos.NewSQLMarketsRepository(cfg)
	market := marketsRepo.GetMarketByName("mobile.de")
	propertiesRepository := repos.NewPropertiesRepository(cfg)
	return &URLBuilder{
		modelsMap:              initModelsAdapterMap(),
		brandModelParamsValues: initParamNames(),
		fuelParam:              initFuelParams(),
		propertiesRepo:         repos.NewPropertiesRepository(cfg),
		marketsRepo:            marketsRepo,
		marketID:               market.ID,
		propertiesRepository:   propertiesRepository,
	}
}

func (b URLBuilder) GetPageURL(criteria jobs.Criteria, pageNumber int) string {
	brand := criteria.Brand
	model := b.modelsMap[criteria.CarModel]
	brandParamValue := b.brandModelParamsValues[brand][model].Brand
	modelParamValue := b.brandModelParamsValues[brand][model].Model
	subModelParamValue := ""
	if b.brandModelParamsValues[brand][model].SubModel != nil {
		subModelParamValue = *b.brandModelParamsValues[brand][model].SubModel
	}

	fuelParam := ""
	fuel := b.fuelParam[criteria.Fuel]
	if fuel != "" {
		fuelParam = fmt.Sprintf(",ful:%s", fuel)
	}

	url := fmt.Sprintf("https://www.mobile.de/ro/automobil/%s-%s/vhc:car,pgn:%d,pgs:50,srt:price,sro:asc,ms1:%s_%s_%s,frn:%d%s,mlx:125000,dmg:false", brand, model, pageNumber, brandParamValue, modelParamValue, subModelParamValue, *criteria.YearFrom, fuelParam)

	return url
}

type MobileDEURLMap struct {
	PathBrandName      string
	PathModelName      string
	PageNumber         int
	QueryParamBrand    string
	QueryParamModel    string
	QueryParamSubModel string
	QueryParamYearFrom *int
	QueryParamFuel     string
}

func (b URLBuilder) GetURL(criteria jobs.Criteria, pageNumber int) string {
	brand := b.propertiesRepo.GetPropertyMarketValuesForTypeAndValue(repos.Brand, criteria.Brand, b.marketID).Value
	model := b.propertiesRepo.GetPropertyMarketValuesForTypeAndValue(repos.Model, criteria.CarModel, b.marketID).Value
	queryParamBrand := b.propertiesRepo.GetPropertyMarketValuesForTypeAndValue(repos.QueryParamBrandParam, criteria.Brand, b.marketID).Value
	queryParamModel := b.propertiesRepo.GetPropertyMarketValuesForTypeAndValue(repos.QueryParamModelParam, criteria.CarModel, b.marketID).Value
	queryParamFuel := b.propertiesRepo.GetPropertyMarketValuesForTypeAndValue(repos.QueryParamFuelParam, criteria.Fuel, b.marketID).Value
	//queryParamSubModel := b.propertiesRepo.GetPropertyMarketValuesForTypeAndValue(repos.SubModel, criteria., b.marketID).Value
	// TODO add submodel in criteria ?
	queryParamSubModel := ""

	paramsStruct := MobileDEURLMap{
		PathBrandName:      brand,
		PathModelName:      model,
		PageNumber:         pageNumber,
		QueryParamBrand:    queryParamBrand,
		QueryParamModel:    queryParamModel,
		QueryParamSubModel: queryParamSubModel,
		QueryParamYearFrom: criteria.YearFrom,
		QueryParamFuel:     queryParamFuel,
	}
	//https://suchen.mobile.de/fahrzeuge/search.html?dam=false&fe=HYBRID_PLUGIN&fr=2019%3A&isSearchRequest=true&ml=%3A125000&ms=3500%3B49%3B%3B&ref=dsp&s=Car&sb=rel&vc=Car
	url := fmt.Sprintf("https://www.mobile.de/ro/automobil/%s-%s/vhc:car,pgn:%d,pgs:50,srt:price,sro:asc,ms1:%s_%s_%s,frn:%d%s,mlx:125000,dmg:false", paramsStruct.PathBrandName, paramsStruct.PathModelName, paramsStruct.PageNumber, paramsStruct.QueryParamBrand, paramsStruct.QueryParamModel, paramsStruct.QueryParamSubModel, *criteria.YearFrom, paramsStruct.QueryParamFuel)
	return url
}

func initFuelParams() map[string]string {
	fuelMap := make(map[string]string)
	fuelMap["diesel"] = "diesel"
	fuelMap["petrol"] = "petrol"
	fuelMap["hybrid-petrol"] = "hybrid_petrol"
	fuelMap["hybrid-diesel"] = "hybrid_diesel"
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
	modelsMap["cx-60"] = "cx-60"
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

	mazdaModelsMap := map[string]BrandModelValues{}
	cx60 := BrandModelValues{
		Brand: "16800",
		Model: "63",
	}
	mazdaModelsMap["cx-60"] = cx60
	params["mazda"] = mazdaModelsMap

	return params
}
