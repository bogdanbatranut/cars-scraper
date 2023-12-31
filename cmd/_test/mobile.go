package main

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/markets"

	"github.com/google/uuid"
)

func main() {

	// https://www.mobile.de/ro/automobil/mercedes-benz-clasa-gle/vhc:car,srt:price,sro:asc,ms1:17200_-58_,frn:2019,ful:diesel,mlx:125000,dmg:false

	criteria := adsdb.Criteria{
		Brand:        "mercedes-benz",
		CarModel:     "gle_class",
		YearFrom:     toIntPointer(2019),
		YearTo:       nil,
		Fuel:         "diesel",
		KmFrom:       nil,
		KmTo:         toIntPointer(125000),
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

	implementationStragegies := markets.NewImplemetationStrategies()
	impl := implementationStragegies.GetImplementation("mobile.de")
	impl.Execute(rsj)
}

func toIntPointer(value int) *int {
	return &value
}

func toStringPointer(value string) *string {
	return &value
}
