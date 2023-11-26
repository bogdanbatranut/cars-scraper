package main

import (
	"carscraper/pkg/adsdb"
)

func main() {
	//
	//market1 := adsdb.Market{
	//	Name: "autovit",
	//	PageURL:  "www.autovit.ro",
	//}
	//
	//market2 := adsdb.Market{
	//	Name: "mobile.de",
	//	PageURL:  "www.mobile.de",
	//}
	//
	//var mkts []adsdb.Market
	//mkts = append(mkts, market1)
	//mkts = append(mkts, market2)
	//
	criteria := adsdb.Criteria{
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
	//
	//log.Println(criteria)
	//
	repo := adsdb.NewRepository("carsfinder")

	repo.AddCriteria(criteria)
	//err := repo.Migrate()
	//
	//errorshandler.HandleErr(err)
	//
	//repo.AddCriteria(criteria)

	//databaseName := "carsfinder"
	//dsn := fmt.Sprintf("root:rootpass@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", databaseName)
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//
	//errorshandler.HandleErr(err)
	//
	//jc := jobs.NewJobsCoordinatorService(db)
	//jc.Start()

}

func pOf(n int) *int {
	return &n
}
