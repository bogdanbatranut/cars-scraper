package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/jobs"
	"carscraper/pkg/repos"

	"github.com/google/uuid"
)

func main() {

	// create jobs from criteria
	//cfg, err := amconfig.NewViperConfig()
	//errorshandler.HandleErr(err)

	//adapter := NewCriteriasJobsAdapter(cfg)
	//job := adapter.CreateJob(6, 12)

	// get implementation and execute jobs
	//logger := logging.NewScrapeLoggingService(cfg)
	//collector := autoscoutrodcollector.NewAutoScoutRodCollector(nil, nil)
	//implementation := autoscout.NewAutoscoutStrategy(&logger, collector, nil)
	//ads, b, err := implementation.Execute(*job)
	//if err != nil {
	//	return
	//}
	//log.Println(ads)
	//log.Println(b)
	//log.Println("starting page scraping service...")

	//sjc := temp.NewPageScrapingService(cfg, temp.WithSimpleMessageQueueRepository(cfg))
	//sjc.Start()
}

type CriteriasJobsAdapter struct {
	criteriasRepo *repos.SQLCriteriaRepository
	marketsRepo   *repos.SQLMarketsRepository
}

func NewCriteriasJobsAdapter(config amconfig.IConfig) *CriteriasJobsAdapter {
	return &CriteriasJobsAdapter{
		criteriasRepo: repos.NewSQLCriteriaRepository(config),
		marketsRepo:   repos.NewSQLMarketsRepository(config),
	}
}

func (adapter CriteriasJobsAdapter) CreateJob(criteriaID uint, marketID uint) *jobs.SessionJob {
	criteria := adapter.criteriasRepo.GetCriteriaByID(criteriaID)
	market := adapter.marketsRepo.GetMarketByID(marketID)
	job := jobs.SessionJob{
		SessionID:  uuid.New(),
		JobID:      uuid.New(),
		CriteriaID: criteria.ID,
		MarketID:   marketID,
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
			Name:       market.Name,
			PageNumber: 1,
		},
	}
	return &job
}
