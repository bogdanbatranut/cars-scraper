package main

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/markets/bmw_de"

	"github.com/google/uuid"
)

func main() {
	adapter := bmw_de.NewBMWDECollyMarketAdapter()
	sessionID := uuid.New()
	pageN := 1
	for {
		job := createJob(sessionID, pageN)
		res := adapter.GetAds(job)
		if res.IsLastPage {
			break
		}
		if len(*res.Ads) == 0 {
			break
		}
		pageN++
	}
}

func createJob(sessionID uuid.UUID, pageNumber int) *jobs.SessionJob {
	yearFrom := 2019
	yearTo := 2024
	kmFrom := 0
	kmTo := 125000
	job := jobs.SessionJob{
		AllowIncrementPage: true,
		SessionID:          sessionID,
		JobID:              uuid.New(),
		Criteria: jobs.Criteria{
			Brand:    "bmw",
			CarModel: "7-series",
			YearFrom: &yearFrom,
			YearTo:   &yearTo,
			Fuel:     "diesel",
			KmFrom:   &kmFrom,
			KmTo:     &kmTo,
		},
		Market: jobs.Market{
			Name:       "gebraucht_bmw",
			PageNumber: pageNumber,
		},
	}
	return &job
}
