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
}

func NewURLBuilder(criteria jobs.Criteria) *URLBuilder {
	return &URLBuilder{
		criteria:     criteria,
		paramsMapper: NewParamsMapper(),
	}
}

func (b URLBuilder) GetPageURL(pageNumber int) string {
	variables := b.createVariablesParam(pageNumber)
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

	url := fmt.Sprintf("https://www.autovit.ro/graphql?operationName=listingScreen&variables=%s&extensions=%s", encodedVariablesStr, encodedExtensionsStr)
	return url
}

func (b URLBuilder) createVariablesParam(page int) VariablesParam {
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
		Parameters:                 []string{"make", "vat", "fuel_type", "mileage", "engine_capacity", "engine_code", "engine_power", "first_registration_year", "model", "version", "year"},
		SearchTerms:                nil,
		SortBy:                     "filter_float_price:asc",
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
		{
			Name:  "filter_enum_fuel_type",
			Value: b.criteria.Fuel,
		},
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

	return f
}

func (b URLBuilder) createExtensionsParam() ExtensionsParam {
	pq := PersistedQuery{
		Sha256Hash: "ea42916db1b919c901d17722dc529de452fa5071e8695743fb2d5378a9dc0315",
		Version:    1,
	}
	return ExtensionsParam{
		PersistedQuery: pq,
	}
}
