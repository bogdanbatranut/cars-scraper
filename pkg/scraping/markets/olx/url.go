package olx

import (
	"carscraper/pkg/jobs"
	"fmt"
	url2 "net/url"
)

type BrandModelIds struct {
	BrandID string
	ModelID string
}

type URLBuilder struct {
	criteria        jobs.Criteria
	brandModelsMapp map[string]map[string]*BrandModelIds
}

func NewURLBuilder(criteria jobs.Criteria) *URLBuilder {
	return &URLBuilder{
		criteria:        criteria,
		brandModelsMapp: buildBrandModelIDsParams(),
	}
}

func (b URLBuilder) GetPageURL(pageNumber int) *string {
	// https://www.olx.ro/api/v1/offers/?offset=0&limit=40&category_id=183&currency=EUR&filter_enum_model%5B0%5D=7-es-sorozat&filter_enum_petrol%5B0%5D=diesel&filter_float_rulaj_pana%3Ato=125000&filter_float_year%3Afrom=2019&filter_refiners=spell_checker&suggest_filters=true&sl=18d1b904943xf2f313c
	limit := 40
	offset := (pageNumber - 1) * limit
	kmTo := *b.criteria.KmTo
	fuel := b.criteria.Fuel
	yearFrom := *b.criteria.YearFrom
	if b.brandModelsMapp[b.criteria.Brand][b.criteria.CarModel] == nil {
		return nil
	}
	brand := b.brandModelsMapp[b.criteria.Brand][b.criteria.CarModel].BrandID
	model := b.brandModelsMapp[b.criteria.Brand][b.criteria.CarModel].ModelID

	fem := url2.QueryEscape("filter_enum_model[0]")
	filterEnumModel := fmt.Sprintf("%s=%s", fem, model)

	fep := url2.QueryEscape("filter_enum_petrol[0]")
	filterEnumPetrol := fmt.Sprintf("%s=%s", fep, fuel)

	ffrp := url2.QueryEscape("filter_float_rulaj_pana:to")
	kmToStr := fmt.Sprintf("%s=%d", ffrp, kmTo)

	ffyf := url2.QueryEscape("filter_float_year:from")
	yearFromStr := fmt.Sprintf("%s=%d", ffyf, yearFrom)

	url := fmt.Sprintf("https://www.olx.ro/api/v1/offers/?offset=%d&limit=40&category_id=%s&currency=EUR&%s&%s&%s&%s&filter_refiners=spell_checker&suggest_filters=true&sl=18d1b904943xf2f313c", offset, brand, filterEnumModel, filterEnumPetrol, kmToStr, yearFromStr)
	return &url
}

func buildBrandModelIDsParams() map[string]map[string]*BrandModelIds {
	brandModelsMap := make(map[string]map[string]*BrandModelIds)
	bmw := "183"
	bmwx5 := BrandModelIds{
		BrandID: bmw,
		ModelID: "x5",
	}
	modelMap := make(map[string]*BrandModelIds)
	modelMap["x5"] = &bmwx5
	brandModelsMap["bmw"] = modelMap

	bmwx5m := BrandModelIds{
		BrandID: bmw,
		ModelID: "bmw-x5m",
	}
	modelMap["x5-m"] = &bmwx5m
	brandModelsMap["bmw"] = modelMap

	bmwx4 := BrandModelIds{
		BrandID: bmw,
		ModelID: "x4",
	}
	modelMap["x4"] = &bmwx4
	brandModelsMap["bmw"] = modelMap

	bmwx6 := BrandModelIds{
		BrandID: bmw,
		ModelID: "x6",
	}
	modelMap["x6"] = &bmwx6
	brandModelsMap["bmw"] = modelMap

	bmwx6m := BrandModelIds{
		BrandID: bmw,
		ModelID: "bmw-x6m",
	}
	modelMap["x6-m"] = &bmwx6m
	brandModelsMap["bmw"] = modelMap

	bmw7 := BrandModelIds{
		BrandID: bmw,
		ModelID: "7-es-sorozat",
	}
	modelMap["7-series"] = &bmw7
	brandModelsMap["bmw"] = modelMap

	bmwx3 := BrandModelIds{
		BrandID: bmw,
		ModelID: "x3",
	}
	modelMap["x3"] = &bmwx3
	brandModelsMap["bmw"] = modelMap

	modelMap = make(map[string]*BrandModelIds)
	mb := "195"
	//mbglcCoupe := BrandModelIds{
	//	BrandID: mb,
	//	ModelID: "5e404d08-99ff-444a-a9b1-64101c387488",
	//}
	//modelMap["glc-coupe"] = &mbglcCoupe
	//brandModelsMap["mercedes-benz"] = modelMap

	mbglc := BrandModelIds{
		BrandID: mb,
		ModelID: "glc",
	}
	modelMap["glc-class"] = &mbglc
	brandModelsMap["mercedes-benz"] = modelMap

	mbglb := BrandModelIds{
		BrandID: mb,
		ModelID: "glb",
	}
	modelMap["glb-class"] = &mbglb
	brandModelsMap["mercedes-benz"] = modelMap

	mbgle := BrandModelIds{
		BrandID: mb,
		ModelID: "gle",
	}
	modelMap["gle-class"] = &mbgle

	mbe := BrandModelIds{
		BrandID: mb,
		ModelID: "e-class",
	}
	modelMap["e-class"] = &mbe

	brandModelsMap["mercedes-benz"] = modelMap

	modelMap = make(map[string]*BrandModelIds)
	opel := "198"
	mokka := BrandModelIds{
		BrandID: opel,
		ModelID: "mokka",
	}
	modelMap["mokka"] = &mokka
	brandModelsMap["opel"] = modelMap

	modelMap = make(map[string]*BrandModelIds)
	skoda := "203"
	octavia := BrandModelIds{
		BrandID: skoda,
		ModelID: "octavia",
	}
	modelMap["octavia"] = &octavia
	brandModelsMap["skoda"] = modelMap

	superb := BrandModelIds{
		BrandID: skoda,
		ModelID: "superb",
	}
	modelMap["superb"] = &superb
	brandModelsMap["skoda"] = modelMap

	modelMap = make(map[string]*BrandModelIds)
	volvo := "208"
	xc90 := BrandModelIds{
		BrandID: volvo,
		ModelID: "xc-90",
	}
	modelMap["xc90"] = &xc90
	brandModelsMap["volvo"] = modelMap

	xc60 := BrandModelIds{
		BrandID: volvo,
		ModelID: "xc-60",
	}
	modelMap["xc60"] = &xc60
	brandModelsMap["volvo"] = modelMap

	s90 := BrandModelIds{
		BrandID: volvo,
		ModelID: "s90",
	}
	modelMap["s90"] = &s90
	brandModelsMap["volvo"] = modelMap

	return brandModelsMap
}
