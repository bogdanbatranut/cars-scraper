package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/properties"
	"carscraper/pkg/repos"
	"log"
)

func testSelects() {

}

func main() {

	cfg, err := amconfig.NewViperConfig()
	if err != nil {
		panic(err)
	}

	propRepo := repos.NewPropertiesRepository(cfg)
	propRepo.Migrate()

	//mkp := propRepo.GetPropertyMarketValues(12, 2)
	////propRepo.GetMarketProperties()
	//log.Println(mkp)

	mkp := propRepo.GetPropertyMarketValuesForTypeAndValue("brand", "bmw", 2)

	log.Println(mkp)

	automallProperty := properties.AutoMallPropertyKey{
		Name:  "brand",
		Value: "bmw",
	}

	marketProperties := []properties.MarketProperty{
		{
			Value:               "inserted market 1 prop 122",
			MarketID:            1,
			AutoMallPropertyKey: automallProperty,
		},
		{
			Value:               "inserted market 2 prop 1",
			MarketID:            2,
			AutoMallPropertyKey: automallProperty,
		},
	}
	automallProperty.MarketProperties = marketProperties

	propRepo.Test(automallProperty)
	propRepo.AddMarketProperties(automallProperty)

	log.Println(propRepo.GetPropertyMarketValues(13, 1))
}
