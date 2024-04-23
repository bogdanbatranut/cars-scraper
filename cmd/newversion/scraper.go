package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/jobs"
	"carscraper/pkg/repos"
	"carscraper/pkg/scraping/scrapingservices"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
)

func main() {

	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)
	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		log.Println("canceling")
		cancel()
		done <- true
	}()

	fmt.Println("awaiting signal")

	//runner := scrapingservices.NewServicesRunner(ctx, cfg)
	//handler := runner.StartServices()

	scrapingMapper := scrapingservices.NewScrapingAdaptersMapper()

	rodScrapingService := scrapingservices.NewRodScrapingService(ctx, scrapingMapper, cfg)
	rodScrapingService.Start()

	collyScrapingService := scrapingservices.NewCollyScrapingService(ctx, scrapingMapper)
	collyScrapingService.Start()

	jsonScrapingService := scrapingservices.NewJSONScrapingService(ctx, scrapingMapper)
	jsonScrapingService.Start()

	//sjh := scrapingservices.NewSessionJobHandler(ctx, cfg, rodScrapingService, collyScrapingService, jsonScrapingService)
	sjh := scrapingservices.NewSessionJobHandler(ctx, cfg,
		scrapingservices.WithMarketService("autovit", jsonScrapingService),
		scrapingservices.WithMarketService("mobile.de", collyScrapingService),
		scrapingservices.WithMarketService("autoscout", rodScrapingService),
		scrapingservices.WithMarketService("autotracknl", rodScrapingService),
		scrapingservices.WithMarketService("olx", jsonScrapingService),
	)

	sjh.StartWithoutMQ()
	//
	adapter := NewCriteriasJobsAdapter(cfg)
	go func() {
		markets := []uint{9, 11, 12, 13}
		criterias := []uint{20, 21}
		//marketID := uint(9)
		for _, marketID := range markets {
			if marketID == 10 {
				marketID++
				continue
			}
			for _, criteriaID := range criterias {
				// criteria 7 volvo s90
				// criteria 6 bmw 7 series
				job := adapter.CreateJob(uuid.New(), criteriaID, marketID)
				sjh.AddScrapingJob(*job)
				//job = adapter.CreateJob(uuid.New(), 8, marketID)
				//sjh.AddScrapingJob(*job)
			}
		}
		done <- true
	}()

	<-done
	fmt.Println("exiting")

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

func (adapter CriteriasJobsAdapter) CreateJob(sessionID uuid.UUID, criteriaID uint, marketID uint) *jobs.SessionJob {
	criteria := adapter.criteriasRepo.GetCriteriaByID(criteriaID)
	market := adapter.marketsRepo.GetMarketByID(marketID)
	job := jobs.SessionJob{
		SessionID:  sessionID,
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
