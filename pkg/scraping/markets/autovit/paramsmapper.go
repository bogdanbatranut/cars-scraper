package autovit

// params mapper

import (
	"fmt"
)

// https://www.mobile.de/
//ro/automobil/mercedes-benz-clasa-gle/vhc:car,pgn:1,pgs:50,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000,dmg:false
// url builder

type ParamsMapper struct {
	modelsMap map[string]string
}

func NewParamsMapper() ParamsMapper {
	return ParamsMapper{
		modelsMap: initModelsAdapterMap(),
	}
}

func (pm ParamsMapper) GetModelParamValue(model string) string {
	return pm.modelsMap[model]
}

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

func initModelsAdapterMap() map[string]string {
	modelsMap := make(map[string]string)

	modelsMap["x4"] = "x4"
	modelsMap["x4-m"] = "x4"
	modelsMap["x5"] = "x5"
	modelsMap["x5-m"] = "bmw-x5m"
	modelsMap["x6"] = "x6"
	modelsMap["x6-m"] = "x6"
	modelsMap["7-series"] = "seria-7"
	modelsMap["gle-class"] = "gle"
	modelsMap["e-class"] = "e"
	modelsMap["s90"] = "s90"
	modelsMap["xc90"] = "xc90"
	modelsMap["xc60"] = "xc60"

	return modelsMap
}

func initParamNames() map[string]map[string]BrandModelValues {
	params := make(map[string]map[string]BrandModelValues)

	mbModelsMap := map[string]BrandModelValues{}
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
