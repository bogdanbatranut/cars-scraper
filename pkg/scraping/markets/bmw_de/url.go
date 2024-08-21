package bmw_de

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/statistics/calculators/helpers"
	"fmt"
	"math"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type URLBuilder struct {
}

func NewURLBuilder() *URLBuilder {
	return &URLBuilder{}
}

func (b URLBuilder) GetPageURL(job jobs.SessionJob) string {
	pageNumber := job.Market.PageNumber
	//return "https://gebrauchtwagen.bmw.de/nsc/search?q=:relevance:condition-firstRegistrationYear:2019:condition-mileageRange:%3C%20150,000%20km::price-asc:series:X::price-asc:model:X4::price-asc:environment-fuelType:Diesel&page=2"
	//        https://gebrauchtwagen.bmw.de/nsc/search?q=:relevance:condition-firstRegistrationYear:1374389604920:condition-mileageRange%3A%3C+%25%21f%28%2Aint%3D0x140000112a0%29+km::price-asc:series:X::price-asc:model:X4::price-asc:environment-fuelType:Diesel&page=1
	yearFrom := *job.Criteria.YearFrom
	kmTo := job.Criteria.KmTo
	series := getCarSeriesFromCarModel(job.Criteria.CarModel)
	model := getCarModelFromCriteria(job.Criteria.CarModel)
	fuel := getFuelTypeFromCriteria(job.Criteria.Fuel)

	kmFloat := float64(*kmTo)
	num := float64(kmFloat / 100000)
	num = RRound(num, 0.5)
	kmFloat = num * 100000

	p := message.NewPrinter(language.English)
	kmWithCommaFormated := p.Sprintf("%.0f", kmFloat)

	firstPart := "condition-mileageRange:%3C%20"
	lastPart := "%20km"

	encodedConditionMileageRange := fmt.Sprintf("%s%s%s", firstPart, kmWithCommaFormated, lastPart)

	addModelparam := addModelParam(job.Criteria)
	modelParam := ""
	if addModelparam {
		modelParam = fmt.Sprintf(":model:%s", model)
	}

	url := fmt.Sprintf("https://gebrauchtwagen.bmw.de/nsc/search?q=:relevance:condition-firstRegistrationYear:%d:%s::price-asc:series:%s::price-asc%s::price-asc:environment-fuelType:%s&page=%d", yearFrom, encodedConditionMileageRange, series, modelParam, fuel, pageNumber)
	return url
}

func RRound(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}

func getCarSeriesFromCarModel(carModel string) string {
	if helpers.InStringArray(carModel, []string{"x4", "x3", "x5", "x6"}) {
		return "X"
	}
	if carModel == "7-series" {
		return "7"
	}
	return ""
}

func getCarModelFromCriteria(criteriaCarModel string) string {
	if criteriaCarModel == "x4" {
		return "X4"
	}
	if criteriaCarModel == "x5" {
		return "X4"
	}
	if criteriaCarModel == "x6" {
		return "X6"
	}
	if criteriaCarModel == "x3" {
		return "X3"
	}
	return ""
}

func getFuelTypeFromCriteria(criteriaFuel string) string {
	fuel := ""
	if criteriaFuel == "diesel" {
		return "Diesel"
	}
	if criteriaFuel == "hybrid-petrol" {
		return "Benzin-Hybrid"
	}
	if criteriaFuel == "hybrid-diesel" {
		return "Plug-In+Hybrid"
	}
	if criteriaFuel == "petrol" {
		return "Benzin"
	}
	if criteriaFuel == "electric" {
		return "Elektrisch"
	}
	if criteriaFuel == "phev-diesel" {
		return "Plug-In+Hybrid"
	}
	if criteriaFuel == "phev-diesel" {
		return "Plug-In+Hybrid"
	}
	return fuel
}

func addModelParam(criteria jobs.Criteria) bool {
	if criteria.CarModel == "7-series" {
		return false
	}
	return true
}
