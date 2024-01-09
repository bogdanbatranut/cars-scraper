package mobile

import (
	"carscraper/pkg/jobs"
	"fmt"
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
}

func NewURLBuilder(criteria jobs.Criteria) *URLBuilder {
	return &URLBuilder{
		criteria:               criteria,
		modelsMap:              initModelsAdapterMap(),
		brandModelParamsValues: initParamNames(),
	}
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
	fuel := b.criteria.Fuel
	url := fmt.Sprintf("https://www.mobile.de/ro/automobil/%s-%s/vhc:car,pgn:%d,pgs:50,srt:price,sro:asc,ms1:%s_%s_%s,frn:2019,ful:%s,mlx:125000,dmg:false", brand, model, pageNumber, brandParamValue, modelParamValue, subModelParamValue, fuel)
	return url
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
	modelsMap["glb-class"] = "clasa-glb"
	modelsMap["x3"] = "x3"
	modelsMap["glc-class"] = "clasa-glc"
	modelsMap["glc-coupe"] = "clasa-glc-coupe"
	modelsMap["octavia"] = "octavia"
	modelsMap["superb"] = "superb"
	modelsMap["mokka"] = "mokka"

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

	return params
}
