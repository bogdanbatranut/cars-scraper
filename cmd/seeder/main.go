package main

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/repos"
)

func main() {
	m := getMarkets()
	c := getGriterias()

	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)
	migrator := repos.NewMigrationRepository(cfg)
	markets := migrator.WriteMarkets(*m)
	migrator.WriteCriterias(markets, c)

}

func getMarkets() *[]adsdb.Market {
	market1 := adsdb.Market{
		Name: "autovit",
		URL:  "www.autovit.ro",
	}

	market2 := adsdb.Market{
		Name: "mobile.de",
		URL:  "www.mobile.de",
	}

	var mkts []adsdb.Market
	mkts = append(mkts, market1)
	mkts = append(mkts, market2)

	return &mkts
}

func getGriterias() *[]adsdb.Criteria {
	c1 := &adsdb.Criteria{
		Brand:        "mazda",
		CarModel:     "cx-5",
		YearFrom:     pOf(2019),
		YearTo:       pOf(2023),
		Fuel:         "diesel",
		KmFrom:       pOf(0),
		KmTo:         pOf(125000),
		AllowProcess: true,
		Markets:      nil,
		ScrapeLogs:   nil,
	}
	var criterias []adsdb.Criteria
	criterias = append(criterias, *c1)
	return &criterias
}

func pOf(n int) *int {
	return &n
}
