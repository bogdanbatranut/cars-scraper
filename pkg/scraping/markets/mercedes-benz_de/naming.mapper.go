package mercedes_benz_de

import (
	"carscraper/pkg/jobs"
	"fmt"
	"strconv"
	"time"
)

type MercedesBenzDENamingMapper struct {
	fuelTypes    map[string]CodesTextEntry
	salesClasses map[string]CodesTextEntry
	bodyGroups   map[string]CodesTextEntry
}

func NewMercedesBenzDENamingMapper() *MercedesBenzDENamingMapper {
	return &MercedesBenzDENamingMapper{
		fuelTypes:    initFuels(),
		salesClasses: initSalesClasses(),
		bodyGroups:   initBodyGroups(),
	}
}

func (mapper MercedesBenzDENamingMapper) GetBrand(brand string) *string {
	return nil
}

func (mapper MercedesBenzDENamingMapper) GetModelCode(model string) *string {
	return nil
}

func (mapper MercedesBenzDENamingMapper) GetModelText(model string) *string {
	return nil
}

func (mapper MercedesBenzDENamingMapper) GetFuelCode(fuel string) *string {
	return nil
}

func (mapper MercedesBenzDENamingMapper) GetFuelText(fuel string) *string {
	return nil
}

func (mapper MercedesBenzDENamingMapper) GetEngineTypeCodesTextEntry(fuel string) CodesTextEntry {
	return mapper.fuelTypes[fuel]
}

func (mapper MercedesBenzDENamingMapper) GetSalesClassCodesTextEntry(model string) CodesTextEntry {
	return mapper.fuelTypes[model]
}

func (mapper MercedesBenzDENamingMapper) GetMarketCriteria(criteria jobs.Criteria) Criteria {
	marketCriteria := Criteria{
		EngineType:        []CodesTextEntry{mapper.fuelTypes[criteria.Fuel]},
		FirstRegistration: mapper.getRegistrationInfo(criteria),
		SalesClass:        []CodesTextEntry{mapper.salesClasses[criteria.CarModel]},
		BodyGroup:         []CodesTextEntry{mapper.bodyGroups[criteria.CarModel]},
	}
	return marketCriteria
}

func (mapper MercedesBenzDENamingMapper) getRegistrationInfo(criteria jobs.Criteria) FirstRegistration {

	yearFrom := *criteria.YearFrom
	yearFrom = yearFrom * 10000
	yearFrom = yearFrom + 101
	year, month, day := time.Now().Date()
	yearToStr := fmt.Sprintf("%d%02d%02d", year, month, day)
	yearTo, err := strconv.Atoi(yearToStr)
	if err != nil {
		panic(err)
	}

	return FirstRegistration{
		Max: yearTo,
		Min: yearFrom,
	}

}

func initFuels() map[string]CodesTextEntry {
	fuels := make(map[string]CodesTextEntry)

	fuels["diesel"] = CodesTextEntry{
		Codes: []string{"1"},
		Text:  "Diesel",
	}

	fuels["petrol"] = CodesTextEntry{
		Codes: []string{"2"},
		Text:  "Benzina",
	}

	fuels["hybrid"] = CodesTextEntry{
		Codes: []string{"6"},
		Text:  "Hibrid",
	}

	fuels["electric"] = CodesTextEntry{
		Codes: []string{"7"},
		Text:  "Electric",
	}
	return fuels
}

func initSalesClasses() map[string]CodesTextEntry {
	salesClasses := make(map[string]CodesTextEntry)

	salesClasses["gle-class"] = CodesTextEntry{
		Codes: []string{"GLE"},
		Text:  "GLE",
	}

	salesClasses["e-class"] = CodesTextEntry{
		Codes: []string{"E"},
		Text:  "Clasa E",
	}

	salesClasses["glc-class"] = CodesTextEntry{
		Codes: []string{"GLC"},
		Text:  "GLC",
	}

	salesClasses["glc-coupe"] = CodesTextEntry{
		Codes: []string{"GLC"},
		Text:  "GLC",
	}

	salesClasses["glb-class"] = CodesTextEntry{
		Codes: []string{"GLB"},
		Text:  "GLB",
	}
	return salesClasses
}

func initBodyGroups() map[string]CodesTextEntry {
	bodyGroups := make(map[string]CodesTextEntry)

	bodyGroups["glc-coupe"] = CodesTextEntry{
		Codes: []string{"3"},
		Text:  "Coupé",
	}

	bodyGroups["glc-class"] = CodesTextEntry{
		Codes: []string{"6"},
		Text:  "SUV &#x2F; Off-Roader",
	}

	bodyGroups["gle-coupe"] = CodesTextEntry{
		Codes: []string{"3"},
		Text:  "Coupé",
	}

	bodyGroups["gle-class"] = CodesTextEntry{
		Codes: []string{"6"},
		Text:  "SUV &#x2F; Off-Roader",
	}

	return bodyGroups
}
