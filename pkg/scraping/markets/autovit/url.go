package autovit

import (
	"carscraper/pkg/jobs"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// url builder
type Filters struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Experiment struct {
	Key     string `json:"key"`
	Variant string `json:"variant"`
}

type ExtensionsParam struct {
	PersistedQuery PersistedQuery `json:"persistedQuery"`
}

type PersistedQuery struct {
	Sha256Hash string `json:"sha256Hash"`
	Version    int    `json:"version"`
}

type VariablesParam struct {
	Click2BuyExperimentId      string       `json:"click2BuyExperimentId"`
	Click2BuyExperimentVariant string       `json:"click2BuyExperimentVariant"`
	Experiments                []Experiment `json:"experiments"`
	Filters                    []Filters    `json:"filters"`
	IncludeClick2Buy           bool         `json:"includeClick2Buy"`
	IncludeFiltersCounters     bool         `json:"includeFiltersCounters"`
	IncludePriceEvaluation     bool         `json:"includePriceEvaluation"`
	IncludePromotedAds         bool         `json:"includePromotedAds"`
	IncludeRatings             bool         `json:"includeRatings"`
	IncludeSortOptions         bool         `json:"includeSortOptions"`
	MaxAge                     int          `json:"maxAge"`
	Page                       int          `json:"page"`
	Parameters                 []string     `json:"parameters"`
	SearchTerms                []string     `json:"searchTerms"`
	SortBy                     string       `json:"sortBy"`
}

type URLBuilder struct {
	criteria     jobs.Criteria
	paramsMapper ParamsMapper
	fuelsMap     map[string][]string
}

func NewURLBuilder(criteria jobs.Criteria) *URLBuilder {
	return &URLBuilder{
		criteria:     criteria,
		paramsMapper: NewParamsMapper(),
		fuelsMap:     initFuelsMap(),
	}
}

func initFuelsMap() map[string][]string {
	fuelMap := make(map[string][]string)
	fuelMap["hybrid"] = []string{"hybrid", "plugin-hybrid"}
	fuelMap["hybrid-petrol"] = []string{"hybrid", "plugin-hybrid"}
	fuelMap["diesel"] = []string{"diesel"}
	fuelMap["petrol"] = []string{"petrol"}
	fuelMap["hybrid-diesel"] = []string{"hybrid", "plugin-hybrid"}
	return fuelMap
}

func (b URLBuilder) GetPageURL(pageNumber int) string {
	//variables := b.createVariablesParam(pageNumber)
	variables := b.createNewVariablesParam(pageNumber)
	variablesBts, err := json.Marshal(variables)
	if err != nil {
		panic(err)
	}
	variablesStr := string(variablesBts)
	extensions := b.createExtensionsParam()
	extensionsBts, err := json.Marshal(extensions)
	if err != nil {
		panic(err)
	}
	extensionsStr := string(extensionsBts)
	encodedVariablesStr := url.QueryEscape(variablesStr)
	encodedExtensionsStr := url.QueryEscape(extensionsStr)

	pageURL := fmt.Sprintf("https://www.autovit.ro/graphql?operationName=listingScreen&variables=%s&extensions=%s", encodedVariablesStr, encodedExtensionsStr)
	return pageURL

	// https://www.autovit.ro/graphql?operationName=listingScreen&variables=%7B%22after%22%3Anull%2C%22experiments%22%3A%5B%7B%22key%22%3A%22MCTA-1414%22%2C%22variant%22%3A%22a%22%7D%2C%7B%22key%22%3A%22MCTA-1617%22%2C%22variant%22%3A%22b%22%7D%2C%7B%22key%22%3A%22MCTA-1660%22%2C%22variant%22%3A%22a%22%7D%2C%7B%22key%22%3A%22MCTA-1661%22%2C%22variant%22%3A%22a%22%7D%2C%7B%22key%22%3A%22CARS-62302%22%2C%22variant%22%3A%22a%22%7D%2C%7B%22key%22%3A%22MCTA-1736%22%2C%22variant%22%3A%22a%22%7D%2C%7B%22key%22%3A%22MCTA-1715%22%2C%22variant%22%3A%22b%22%7D%2C%7B%22key%22%3A%22MCTA-1721%22%2C%22variant%22%3A%22a%22%7D%2C%7B%22key%22%3A%22CARS-64661%22%2C%22variant%22%3A%22a%22%7D%5D%2C%22filters%22%3A%5B%7B%22name%22%3A%22filter_enum_make%22%2C%22value%22%3A%22mercedes-benz%22%7D%2C%7B%22name%22%3A%22filter_enum_model%22%2C%22value%22%3A%22gle%22%7D%2C%7B%22name%22%3A%22filter_enum_fuel_type%22%2C%22value%22%3A%22diesel%22%7D%2C%7B%22name%22%3A%22filter_float_mileage %3Ato%22%2C%22value%22%3A%22125000%22%7D%2C%7B%22name%22%3A%22filter_float_year%3Afrom%22%2C%22value%22%3A%222020%22%7D%2C%7B%22name%22%3A%22order%22%2C%22value%22%3A%22relevance_web%22%7D%2C%7B%22name%22%3A%22category_id%22%2C%22value%22%3A%2229%22%7D%5D%2C%22includeCepik%22%3Afalse%2C%22includeFiltersCounters%22%3Afalse%2C%22includeNewPromotedAds%22%3Afalse%2C%22includePriceEvaluation%22%3Atrue%2C%22includePromotedAds%22%3Afalse%2C%22includeRatings%22%3Afalse%2C%22includeSortOptions%22%3Afalse%2C%22includeSuggestedFilters%22%3Afalse%2C%22maxAge%22%3A60%2C%22page%22%3A1%2C%22parameters%22%3A%5B%22make%22%2C%22vat%22%2C%22fuel_type%22%2C%22mileage%22%2C%22engine_capacity%22%2C%22engine_code%22%2C%22engine_power%22%2C%22first_registration_year%22%2C%22model%22%2C%22version%22%2C%22year%22%5D%2C%22promotedInput%22%3A%7B%7D%2C%22searchTerms%22%3Anull%2C%22sortBy%22%3A%22relevance_web%22%7D&extensions=%7B%22persistedQuery%22%3A%7B%22sha256Hash%22%3A%221a840f0ab7fbe2543d0d6921f6c963de8341e04a4548fd1733b4a771392f900a%22%2C%22version%22%3A1%7D%7D
	// https://www.autovit.ro/graphql?operationName=listingScreen&variables=%7B%22after%22%3Anull%2C%22experiments%22%3A%5B%7B%22key%22%3A%22MCTA-1414%22%2C%22variant%22%3A%22a%22%7D%2C%7B%22key%22%3A%22MCTA-1617%22%2C%22variant%22%3A%22b%22%7D%2C%7B%22key%22%3A%22MCTA-1660%22%2C%22variant%22%3A%22a%22%7D%2C%7B%22key%22%3A%22MCTA-1661%22%2C%22variant%22%3A%22a%22%7D%2C%7B%22key%22%3A%22CARS-62302%22%2C%22variant%22%3A%22a%22%7D%2C%7B%22key%22%3A%22MCTA-1736%22%2C%22variant%22%3A%22a%22%7D%2C%7B%22key%22%3A%22MCTA-1715%22%2C%22variant%22%3A%22b%22%7D%2C%7B%22key%22%3A%22MCTA-1721%22%2C%22variant%22%3A%22b%22%7D%2C%7B%22key%22%3A%22CARS-64661%22%2C%22variant%22%3A%22b%22%7D%5D%2C%22filters%22%3A%5B%7B%22name%22%3A%22filter_enum_make%22%2C%22value%22%3A%22volvo        %22%7D%2C%7B%22name%22%3A%22filter_enum_model%22%2C%22value%22%3A%22xc0%22%7D%2C%7B%22name%22%3A%22filter_float_year%3Afrom%22%2C%22value%22%3A%222019%22%7D%2C%7B%22name%22%3A%22filter_float_mileage%3Ato%22%2C%22value%22%3A%22125000%22%7D%2C%7B%22name%22%3A%22category_id%22%2C%22value%22%3A%2229%22%7D%2C%7B%22name%22%3A%22filter_enum_fuel_type%22%2C%22value%22%3A%22hybrid%22%7D%2C%7B%22name%22%3A%22filter_enum_fuel_type%22%2C%22value%22%3A%22plugin-hybrid%22%7D%5D%2C%22includeCepik%22%3Afalse%2C%22includeFiltersCounters%22%3Afalse%2C%22includeNewPromotedAds%22%3Afalse%2C%22includePriceEvaluation%22%3Atrue%2C%22includePromotedAds%22%3Afalse%2C%22includeRatings%22%3Afalse%2C%22includeSortOptions%22%3Afalse%2C%22includeSuggestedFilters%22%3Afalse%2C%22maxAge%22%3A0%2C%22page%22%3A1%2C%22parameters%22%3A%5B%22make%22%2C%22vat%22%2C%22mileage%22%2C%22engine_capacity%22%2C%22engine_code%22%2C%22engine_power%22%2C%22first_registration_year%22%2C%22model%22%2C%22version%22%2C%22year%22%5D%2C%22promotedInput%22%3A%7B%7D%2C%22searchTerms%22%3Anull%2C%22sortBy%22%3A%22%22%7D&extensions=%7B%22persistedQuery%22%3A%7B%22sha256Hash%22%3A%22ea42916db1b919c901d17722dc529de452fa5071e8695743fb2d5378a9dc0315%22%2C%22version%22%3A1%7D%7D
}

func (b URLBuilder) createNewVariablesParam(page int) VariablesRequestParamValue {
	experiments := createExperiments()
	parameters := []string{"make", "vat", "fuel_type", "mileage", "engine_capacity", "engine_code", "engine_power", "first_registration_year", "model", "version", "year"}

	return VariablesRequestParamValue{
		After:                   nil,
		Experiments:             experiments,
		Filters:                 b.createFiltersFromCriteria(),
		IncludeCepik:            false,
		IncludeFiltersCounters:  false,
		IncludeNewPromotedAds:   false,
		IncludePriceEvaluation:  true,
		IncludePromotedAds:      false,
		IncludeRatings:          false,
		IncludeSortOptions:      false,
		IncludeSuggestedFilters: false,
		MaxAge:                  0,
		Page:                    page,
		Parameters:              parameters,
		PromotedInput:           struct{}{},
		SearchTerms:             nil,
		SortBy:                  "filter_float_price:asc",
	}
}

func createExperiments() []Experiment {
	return []Experiment{
		{Key: "MCTA-1414", Variant: "a"},
		{Key: "MCTA-1617", Variant: "b"},
		{Key: "MCTA-1660", Variant: "a"},
		{Key: "MCTA-1661", Variant: "a"},
		{Key: "CARS-62302", Variant: "a"},
		{Key: "MCTA-1736", Variant: "a"},
		{Key: "MCTA-1715", Variant: "b"},
		{Key: "MCTA-1721", Variant: "b"},
		{Key: "CARS-64661", Variant: "b"},
	}
}

func (b URLBuilder) createVariablesParam(page int) VariablesParam {
	parameters := []string{"make", "vat", "mileage", "engine_capacity", "engine_code", "engine_power", "first_registration_year", "model", "version", "year"}
	if b.criteria.Fuel != "" {
		parameters = append(parameters, "fuel_type")
	}
	return VariablesParam{
		Click2BuyExperimentId:      "",
		Click2BuyExperimentVariant: "",
		Experiments:                b.createExperiments(),
		Filters:                    b.createFiltersFromCriteria(),
		IncludeClick2Buy:           false,
		IncludeFiltersCounters:     false,
		IncludePriceEvaluation:     true,
		IncludePromotedAds:         false,
		IncludeRatings:             false,
		IncludeSortOptions:         false,
		MaxAge:                     60,
		Page:                       page,
		Parameters:                 parameters,
		SearchTerms:                nil,
		SortBy:                     "filter_float_price:asc",
		//Parameters:                 []string{"make", "vat", "fuel_type", "mileage", "engine_capacity", "engine_code", "engine_power", "first_registration_year", "model", "version", "year"},
	}

}

func (b URLBuilder) createExperiments() []Experiment {
	ex := []Experiment{
		{
			Key:     "MCTA-900",
			Variant: "a",
		},
		{
			Key:     "MCTA-1059",
			Variant: "a",
		},
	}
	return ex
}

func (b URLBuilder) createFuelFilters() []Filters {
	var fuelFilters []Filters
	criteriaFuels := b.fuelsMap[b.criteria.Fuel]
	for _, fuelType := range criteriaFuels {
		fuelFilter := Filters{
			Name:  "filter_enum_fuel_type",
			Value: fuelType,
		}
		fuelFilters = append(fuelFilters, fuelFilter)
	}
	return fuelFilters
}

func (b URLBuilder) createFiltersFromCriteria() []Filters {
	f := []Filters{
		{
			Name:  "filter_enum_make",
			Value: b.criteria.Brand,
		},
		{
			Name:  "filter_enum_model",
			Value: b.paramsMapper.GetModelParamValue(b.criteria.CarModel),
		},
		//{
		//	Name:  "filter_enum_fuel_type",
		//	Value: b.criteria.Fuel,
		//},
		{
			Name:  "filter_float_year:from",
			Value: strconv.Itoa(*b.criteria.YearFrom),
		},
		{
			Name:  "filter_float_mileage:to",
			Value: strconv.Itoa(*b.criteria.KmTo),
		},
		{
			Name:  "category_id",
			Value: "29",
		},
	}

	f = append(f, b.createFuelFilters()...)

	//if b.criteria.Fuel != "" {
	//	f = append(f, Filters{
	//		Name:  "filter_enum_fuel_type",
	//		Value: b.criteria.Fuel,
	//	})
	//}

	return f
}

func (b URLBuilder) createExtensionsParam() ExtensionsParam {

	pq := PersistedQuery{
		Sha256Hash: "1a840f0ab7fbe2543d0d6921f6c963de8341e04a4548fd1733b4a771392f900a",
		Version:    1,
	}
	return ExtensionsParam{
		PersistedQuery: pq,
	}
}
