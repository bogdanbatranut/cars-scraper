package main

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/repos"
)

func main() {
	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)

	migrator := repos.NewMigrationRepository(cfg)
	err = migrator.GetDB().AutoMigrate(
		//&adsdb.Criteria{},
		//&adsdb.Ad{},
		//&adsdb.Seller{},
		//&adsdb.Market{},
		//&adsdb.Price{},
		//&adsdb.SellerMarkets{},
		//&adsdb.CriteriaMarkets{},
		//&adsdb.PageLog{},
		//&adsdb.CriteriaLog{},
		&adsdb.SessionLog{},
	)
	errorshandler.HandleErr(err)
}
