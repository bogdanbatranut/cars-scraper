package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/jobs"
	"carscraper/pkg/repos"
	"carscraper/pkg/scraping/markets/olx"
	"log"

	"github.com/google/uuid"
)

func main() {
	implementation := olx.NewOlxStrategy(nil)
	for _, sessionJob := range olxbuildAllSessionJobs() {
		ads, isLastPage, err := implementation.Execute(sessionJob)
		if err != nil {
			panic(err)
		}
		log.Println("IS LAST PAGE ", isLastPage)
		log.Println("AFTER EXEC FOUND: ", len(ads))
	}
	log.Println("DONE")
}

func olxbuildAllSessionJobs() []jobs.SessionJob {
	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)

	repo := repos.NewSQLCriteriaRepository(cfg)

	criterias := repo.GetAll()
	sessionJobs := []jobs.SessionJob{}
	for _, criteria := range *criterias {
		rsj := jobs.SessionJob{
			SessionID:  uuid.New(),
			JobID:      uuid.New(),
			CriteriaID: 1,
			MarketID:   1,
			Criteria: jobs.Criteria{
				Brand:    criteria.Brand,
				CarModel: criteria.CarModel,
				YearFrom: criteria.KmFrom,
				YearTo:   criteria.YearTo,
				Fuel:     criteria.Fuel,
				KmFrom:   criteria.KmFrom,
				KmTo:     criteria.KmTo,
			},
			Market: jobs.Market{
				Name:       "autotrack.nl",
				PageNumber: 1,
			},
		}
		sessionJobs = append(sessionJobs, rsj)
	}
	return sessionJobs
}
