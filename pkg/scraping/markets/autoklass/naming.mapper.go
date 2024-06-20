package autoklass

type AutoklassRONamingMapper struct {
	fuelTypes    map[string]string
	carBrandsIDs map[string]int
	carModelsIDs map[string]int
}

func NewAutoklassRoNamingMapper() *AutoklassRONamingMapper {
	return &AutoklassRONamingMapper{
		fuelTypes:    initFuels(),
		carBrandsIDs: initCarBrandsIDs(),
		carModelsIDs: initCarModelsIDs(),
		//bodyGroups:   initBodyGroups(),
	}
}

func (mapper AutoklassRONamingMapper) GetBrand(brand string) *string {
	return nil
}

func (mapper AutoklassRONamingMapper) GetModelText(model string) *string {
	return nil
}

func (mapper AutoklassRONamingMapper) GetFuelText(fuel string) *string {
	return nil
}

func initFuels() map[string]string {
	fuels := make(map[string]string)
	fuels["diesel"] = "diesel"

	return fuels
}

func initCarBrandsIDs() map[string]int {
	brands := make(map[string]int)
	brands["mercedes-benz"] = 1
	brands["bmw"] = 9
	brands["audi"] = 7
	brands["volkswagen"] = 20
	brands["volvo"] = 18

	return brands
}

func initCarModelsIDs() map[string]int {
	models := make(map[string]int)
	models["gle-class"] = 2
	models["e-class"] = 17
	models["glb-class"] = 19
	models["glc-class"] = 6
	models["glc-coupe"] = 98

	models["7-series"] = 154
	models["x3"] = 111
	models["x5-m"] = 43
	models["x5"] = 114
	models["x6"] = 123

	models["a6"] = 153
	models["q3"] = 350
	models["q5"] = 0
	models["q7"] = 136
	models["q8"] = 0

	models["touareg"] = 202

	models["xc90"] = 167
	models["xc60"] = 351

	return models
}
