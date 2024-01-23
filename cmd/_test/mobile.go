package main

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/jobs"
	"carscraper/pkg/utils"
	"log"

	"github.com/google/uuid"
)

func main() {

	// https://www.mobile.de/ro/automobil/mercedes-benz-clasa-gle/vhc:car,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000,dmg:false

	criteria := adsdb.Criteria{
		Brand:        "mercedes-benz",
		CarModel:     "gle_class",
		YearFrom:     utils.ToIntPointer(2019),
		YearTo:       nil,
		Fuel:         "diesel",
		KmFrom:       nil,
		KmTo:         utils.ToIntPointer(125000),
		AllowProcess: false,
		Markets:      nil,
		ScrapeLogs:   nil,
	}

	rsj := jobs.SessionJob{
		SessionID:  uuid.New(),
		JobID:      uuid.New(),
		CriteriaID: 1,
		MarketID:   1,
		Criteria: jobs.Criteria{
			Brand:    criteria.Brand,
			CarModel: criteria.CarModel,
			YearFrom: criteria.YearFrom,
			YearTo:   criteria.YearTo,
			Fuel:     criteria.Fuel,
			KmFrom:   criteria.KmFrom,
			KmTo:     criteria.KmTo,
		},
		Market: jobs.Market{
			Name:       "mobile.de",
			PageNumber: 1,
		},
	}
	log.Println(rsj)

	//implementationStragegies := markets.NewImplemetationStrategies()
	//impl := implementationStragegies.GetImplementation("mobile.de")
	//impl.Execute(rsj)
}
