package autotrack

import (
	"carscraper/pkg/jobs"
	"fmt"
	"log"
)

type URLBuilder struct {
	brandModels map[string]map[string]*BrandModelIds
}

func NewURLBuilder() *URLBuilder {
	builder := &URLBuilder{
		brandModels: buildBrandModelIDsParams(),
	}
	return builder
}

func (b URLBuilder) GetURL(job jobs.SessionJob) *string {
	fuels := make(map[string]string)
	fuels["diesel"] = "DIESEL"
	fuels["petrol"] = "BENZINE"
	fuels["hybrid-petrol"] = "HYBRIDE_BENZINE"

	bm := b.brandModels[job.Criteria.Brand][job.Criteria.CarModel]
	if bm == nil {
		return nil
	}
	pr := bm.asQueryParams()
	url := fmt.Sprintf(
		"https://www.autotrack.nl/aanbod?minimumbouwjaar=%d&maximumkilometerstand=%d&beschikbaarheidsStatus=beschikbaar&brandstofsoorten=%s&%s&paginanummer=%d&paginagrootte=30&sortering=PRIJS_OPLOPEND",
		*job.Criteria.YearFrom,
		*job.Criteria.KmTo,
		fuels[job.Criteria.Fuel],
		pr,
		job.Market.PageNumber,
	)
	log.Println("AUTOTRACK Fuel : ", fuels[job.Criteria.Fuel], "Fuel in criteria : ", job.Criteria.Fuel)
	return &url
}

type BrandModelIds struct {
	BrandID string
	ModelID string
}

func (bm BrandModelIds) asQueryParams() string {
	return fmt.Sprintf("modelIds=%s&merkIds=%s", bm.ModelID, bm.BrandID)
}

func buildBrandModelIDsParams() map[string]map[string]*BrandModelIds {
	brandModelsMap := make(map[string]map[string]*BrandModelIds)
	bmw := "1a67a3d8-178b-43ee-9071-9ae7f19b316a"
	bmwx5 := BrandModelIds{
		BrandID: bmw,
		ModelID: "e43d83ec-00d4-4cfe-915d-231d267f6d02",
	}
	modelMap := make(map[string]*BrandModelIds)
	modelMap["x5"] = &bmwx5
	brandModelsMap["bmw"] = modelMap

	bmwx5m := BrandModelIds{
		BrandID: bmw,
		ModelID: "9cd35d47-1efe-4048-83f2-fa34df60c370",
	}
	modelMap["x5-m"] = &bmwx5m
	brandModelsMap["bmw"] = modelMap

	bmwx4 := BrandModelIds{
		BrandID: bmw,
		ModelID: "e1a45cec-2cc9-4bfc-879b-becad0313021",
	}
	modelMap["x4"] = &bmwx4
	brandModelsMap["bmw"] = modelMap

	bmwx6 := BrandModelIds{
		BrandID: bmw,
		ModelID: "7477bdee-9d29-4a35-b377-1644f19f1d91",
	}
	modelMap["x6"] = &bmwx6
	brandModelsMap["bmw"] = modelMap

	bmwx6m := BrandModelIds{
		BrandID: bmw,
		ModelID: "75c5d34d-f116-4357-bcb0-782a44e827cb",
	}
	modelMap["x6-m"] = &bmwx6m
	brandModelsMap["bmw"] = modelMap

	bmw7 := BrandModelIds{
		BrandID: bmw,
		ModelID: "9869681e-c3aa-4270-b46e-87ff68c17a3c",
	}
	modelMap["7-series"] = &bmw7
	brandModelsMap["bmw"] = modelMap

	bmwx3 := BrandModelIds{
		BrandID: bmw,
		ModelID: "fbe3e05d-456c-48d2-8e4e-261be8985230",
	}
	modelMap["x3"] = &bmwx3
	brandModelsMap["bmw"] = modelMap

	modelMap = make(map[string]*BrandModelIds)
	mb := "5d1b2626-a08a-464f-aa0e-14a1f80be441"
	mbglcCoupe := BrandModelIds{
		BrandID: mb,
		ModelID: "5e404d08-99ff-444a-a9b1-64101c387488",
	}
	modelMap["glc-coupe"] = &mbglcCoupe
	brandModelsMap["mercedes-benz"] = modelMap

	mbglc := BrandModelIds{
		BrandID: mb,
		ModelID: "3d928ef5-7a53-4cd0-be7c-771bf8b11806",
	}
	modelMap["glc-class"] = &mbglc
	brandModelsMap["mercedes-benz"] = modelMap

	mbglb := BrandModelIds{
		BrandID: mb,
		ModelID: "5e6a1626-26e0-4c23-b28d-b1fbf32e52ae",
	}
	modelMap["glb-class"] = &mbglb
	brandModelsMap["mercedes-benz"] = modelMap

	mbgle := BrandModelIds{
		BrandID: mb,
		ModelID: "8bd9a5ae-fc10-450e-99c1-84ef19dac26c",
	}
	modelMap["gle-class"] = &mbgle

	mbe := BrandModelIds{
		BrandID: mb,
		ModelID: "4be2547b-c266-428d-8fe1-cad32e43a680",
	}
	modelMap["e-class"] = &mbe

	brandModelsMap["mercedes-benz"] = modelMap

	modelMap = make(map[string]*BrandModelIds)
	opel := "7ccf5430-eafb-4042-82c0-43ce39ba1b02"
	mokka := BrandModelIds{
		BrandID: opel,
		ModelID: "85e7360a-cee0-4ae0-85e0-0b595df99471",
	}
	modelMap["mokka"] = &mokka
	brandModelsMap["opel"] = modelMap

	modelMap = make(map[string]*BrandModelIds)
	skoda := "01d8c095-cec4-4904-9001-26115a746977"
	octavia := BrandModelIds{
		BrandID: skoda,
		ModelID: "d9794f12-5b29-4673-9818-c0ef1300a649",
	}
	modelMap["octavia"] = &octavia
	brandModelsMap["skoda"] = modelMap

	superb := BrandModelIds{
		BrandID: skoda,
		ModelID: "a35d0811-8556-471b-ba87-3024cac620ce",
	}
	modelMap["superb"] = &superb
	brandModelsMap["skoda"] = modelMap

	modelMap = make(map[string]*BrandModelIds)
	volvo := "24f2f778-840d-4537-a9d8-119db104ca2e"
	xc90 := BrandModelIds{
		BrandID: volvo,
		ModelID: "a991a613-c209-41a7-8d6e-9ca514b5d99b",
	}
	modelMap["xc90"] = &xc90
	brandModelsMap["volvo"] = modelMap

	xc60 := BrandModelIds{
		BrandID: volvo,
		ModelID: "10683ece-3f3b-423c-b08c-cb99301aed43",
	}
	modelMap["xc60"] = &xc60
	brandModelsMap["volvo"] = modelMap

	s90 := BrandModelIds{
		BrandID: volvo,
		ModelID: "65de8176-8216-4217-93c3-c9e9bafd73b8",
	}
	modelMap["s90"] = &s90
	brandModelsMap["volvo"] = modelMap

	xc40 := BrandModelIds{
		BrandID: volvo,
		ModelID: "8c4dc7aa-755b-4377-ab2f-0d984b5db139",
	}
	modelMap["xc40"] = &xc40
	brandModelsMap["volvo"] = modelMap

	modelMap = make(map[string]*BrandModelIds)

	yarisCross := BrandModelIds{
		BrandID: "adc46765-6df3-4825-bdd2-0c6d427eaa41",
		ModelID: "bf2fdb16-bf85-4d19-9739-92ba48102aa0",
	}

	modelMap["yaris-cross"] = &yarisCross
	brandModelsMap["toyota"] = modelMap

	modelMap = make(map[string]*BrandModelIds)
	vwTouareg := BrandModelIds{
		BrandID: "008d8fc3-b882-4657-9229-4239d5f7e469",
		ModelID: "f74be334-7f97-496e-9bc6-2c2111237461",
	}
	modelMap["touareg"] = &vwTouareg
	brandModelsMap["volkswagen"] = modelMap

	modelMap = make(map[string]*BrandModelIds)
	audiQ7 := BrandModelIds{
		BrandID: "9fdc7e2d-b40c-4ee6-896d-9f0453cb39c6",
		ModelID: "392e2478-d2bb-499c-8fb2-5a5d5cbb1581",
	}
	modelMap["q7"] = &audiQ7
	brandModelsMap["audi"] = modelMap

	audiQ5 := BrandModelIds{
		BrandID: "9fdc7e2d-b40c-4ee6-896d-9f0453cb39c6",
		ModelID: "59ca586c-75ca-4554-94a3-e66be9c6f475",
	}
	modelMap["q5"] = &audiQ5
	brandModelsMap["audi"] = modelMap

	audiA6 := BrandModelIds{
		BrandID: "9fdc7e2d-b40c-4ee6-896d-9f0453cb39c6",
		ModelID: "871154e5-19a7-4190-90ee-89a7a8e00fca",
	}
	modelMap["a6"] = &audiA6
	brandModelsMap["audi"] = modelMap

	audiQ3 := BrandModelIds{
		BrandID: "9fdc7e2d-b40c-4ee6-896d-9f0453cb39c6",
		ModelID: "43aea38c-b114-4469-9b5b-0de166864247",
	}
	modelMap["q3"] = &audiQ3
	brandModelsMap["audi"] = modelMap

	audiQ8 := BrandModelIds{
		BrandID: "9fdc7e2d-b40c-4ee6-896d-9f0453cb39c6",
		ModelID: "c3df76bd-8e5b-40e6-a187-15074ff92524",
	}
	modelMap["q8"] = &audiQ8
	brandModelsMap["audi"] = modelMap

	modelMap = make(map[string]*BrandModelIds)
	mazdacx60 := BrandModelIds{
		BrandID: "112c997b-b9b4-4369-ad2f-c1238d709134",
		ModelID: "8c5d3470-7895-41cd-bed8-bd08e62506d6",
	}
	modelMap["cx-60"] = &mazdacx60
	brandModelsMap["mazda"] = modelMap

	return brandModelsMap
}
