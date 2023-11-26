package main

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/config"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/repos"
)

func main() {
	cfg, err := config.NewViperConfig()
	errorshandler.HandleErr(err)
	migrator := repos.NewMigrationRepository(cfg)
	err = migrator.GetDB().AutoMigrate(
		&adsdb.Criteria{},
		&adsdb.Ad{},
		&adsdb.Seller{},
		&adsdb.Market{},
		&adsdb.Price{},
		&adsdb.SellerMarkets{},
		&adsdb.CriteriaMarkets{},
	)
	errorshandler.HandleErr(err)
}
