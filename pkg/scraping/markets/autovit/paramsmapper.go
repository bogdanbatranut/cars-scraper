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

	modelsMap["x3"] = "x3"
	modelsMap["x3-m"] = "bmw-x3m"
	modelsMap["x4"] = "x4"
	modelsMap["x4-m"] = "x4"
	modelsMap["x5"] = "x5"
	modelsMap["x5-m"] = "bmw-x5m"
	modelsMap["x6"] = "x6"
	modelsMap["x6-m"] = "bmw-x6m"
	modelsMap["7-series"] = "seria-7"
	modelsMap["gle-class"] = "gle"
	modelsMap["e-class"] = "e"
	modelsMap["xc40"] = "xc-40"
	modelsMap["s90"] = "s90"
	modelsMap["xc90"] = "xc-90"
	modelsMap["xc60"] = "xc-60"
	modelsMap["glb-class"] = "glb"
	modelsMap["glc-class"] = "glc"
	modelsMap["glc-coupe"] = "glc-coupe"
	modelsMap["octavia"] = "octavia"
	modelsMap["superb"] = "superb"
	modelsMap["mokka"] = "mokka"
	modelsMap["yaris-cross"] = "mokka"
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
