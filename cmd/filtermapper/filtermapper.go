package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/properties"
	"carscraper/pkg/repos"
	"log"

	"gorm.io/gorm"
)

func createTables(db *gorm.DB) {
	err := db.AutoMigrate(&properties.AutoMallProperty{}, &properties.Market{}, &properties.MarketProperty{})
	if err != nil {
		panic(err)
	}
}

func testSelects() {

}

func main() {

	cfg, err := amconfig.NewViperConfig()
	if err != nil {
		panic(err)
	}

	propRepo := repos.NewPropertiesRepository(cfg)
	ap := propRepo.GetAutoMallPropertyByID(11)
	//propRepo.GetMarketProperties()
	log.Println(ap)
}
