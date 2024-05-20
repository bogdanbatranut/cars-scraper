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
	sjh.Start()

	//https: //www.autotrack.nl/aanbod?minimumbouwjaar=2019&maximumkilometerstand=125000&brandstofsoorten=BENZINE&merkIds=7ccf5430-eafb-4042-82c0-43ce39ba1b02&modelIds=85e7360a-cee0-4ae0-85e0-0b595df99471&beschikbaarheidsStatus=beschikbaar&paginanummer=1&paginagrootte=30&sortering=PRIJS_OPLOPEND
	//https://www.autotrack.nl/aanbod?minimumbouwjaar=2019&maximumkilometerstand=125000&brandstofsoorten=BENZINE&modelIds=85e7360a-cee0-4ae0-85e0-0b595df99471&merkIds=7ccf5430-eafb-4042-82c0-43ce39ba1b02&                                    paginanummer=6&paginagrootte=30&sortering=PRIJS_OPLOPEND	sjh.StartWithoutMQ()
	//
	adapter := NewCriteriasJobsAdapter(cfg)
	go func() {
		//9,2023-11-20 00:06:39.350,2024-03-12 19:57:43.954,,autovit,www.autovit.ro,1
		//10,2023-11-20 00:06:39.350,2024-03-12 19:57:43.963,,webcar,www.webcar.ro,0
		//11,2023-11-20 00:06:39.350,2024-03-12 19:57:43.972,,mobile.de,www.mobile.de,1
		//12,2023-11-20 00:06:39.350,2024-03-12 19:57:43.978,,autoscout,www.autoscout24.ro,1
		//13,2023-11-20 00:06:39.350,2024-03-12 19:57:43.987,,autotracknl,www.autotrack.nl,0
		//14,2023-11-20 00:06:39.350,2024-03-12 19:57:43.996,,olx,www.olx.ro,1

		//markets := []uint{9, 11, 12, 13, 14}
		markets := []uint{11}
		//criterias := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
		criterias := []uint{27}
		//criterias := []uint{22}
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
				//sjh.AddScrapingJobToScrapingService(*job)
				sjh.AddScrapingJob(*job)
				//job = adapter.CreateJob(uuid.New(), 8, marketID)
				//sjh.AddScrapingJobToScrapingService(*job)
			}
		}
		//done <- true
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
