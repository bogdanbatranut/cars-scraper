package mobile

import _const "carscraper/pkg/const"

// params mapper

type MobileCarBrandValues struct {
	MobileCarBrand      string
	MobileCarBrandValue int
	Models              []MobileBrandModels
}

type MobileBrandModels struct {
	Model string
	Value int
}

func CreateMobileCarBrandValues() map[string]MobileCarBrandValues {
	var values map[string]MobileCarBrandValues
	values = make(map[string]MobileCarBrandValues)

	values[_const.MercedesBenz] = MobileCarBrandValues{
		MobileCarBrand:      "Merces-Benz",
		MobileCarBrandValue: 17200,
	}

	return values

}
