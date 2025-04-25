package oferte_bmw

import "carscraper/pkg/jobs"

type RequestPayload struct {
	BrandSeriesVariantIds BrandSeriesVariantIds
	FuelsMap              map[string]int
}

type BrandSeriesVariantIds struct {
	Brands   map[string]int
	Series   map[string]int
	Variants map[string]int
}

type SearchIds struct {
	Brand   int
	Series  int
	Variant int
}

func NewRequestPayload() *RequestPayload {
	return &RequestPayload{
		BrandSeriesVariantIds: initBrandSeriesVariants(),
		FuelsMap:              initFuelsMap(),
	}
}

func (p RequestPayload) GetIds(criteria jobs.Criteria) SearchIds {
	brand := p.BrandSeriesVariantIds.Brands[criteria.Brand]
	s := criteria.CarModel
	series := p.BrandSeriesVariantIds.Series[string(s[0])]
	variant := p.BrandSeriesVariantIds.Variants[criteria.CarModel]

	return SearchIds{
		Brand:   brand,
		Series:  series,
		Variant: variant,
	}
}

func initBrandSeriesVariants() BrandSeriesVariantIds {
	bsvs := BrandSeriesVariantIds{
		Brands:   initBrands(),
		Series:   initSeries(),
		Variants: initVariants(),
	}

	return bsvs
}

func initBrands() map[string]int {
	brands := map[string]int{"bmw": 1}
	return brands
}

func initSeries() map[string]int {
	series := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
		"5": 5,
		"6": 6,
		"7": 7,
		"8": 8,
		"i": 9,
		"m": 20,
		"x": 20,
		"z": 21,
	}

	return series
}

func initVariants() map[string]int {
	variants := map[string]int{
		"x1":       38,
		"x2":       69,
		"x3":       39,
		"x4":       40,
		"x5":       41,
		"x6":       42,
		"x7":       80,
		"7-series": 37,
	}
	return variants
}

func initFuelsMap() map[string]int {
	return map[string]int{
		"petrol":        1,
		"diesel":        2,
		"electric":      3,
		"hybrid-petrol": 4,
		"hybrid-diesel": 6,
	}
}

func (p RequestPayload) GetFuelID(fuel string) int {
	return p.FuelsMap[fuel]
}
