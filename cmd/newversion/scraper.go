package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
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

	loggingService := logging.NewScrapeLoggingService(cfg)

	scrapingMapper := scrapingservices.NewScrapingAdaptersMapper(loggingService)

	rodScrapingService := scrapingservices.NewRodScrapingService(ctx, scrapingMapper, cfg)
	rodScrapingService.Start()

	collyScrapingService := scrapingservices.NewCollyScrapingService(ctx, scrapingMapper)
	collyScrapingService.Start()

	jsonScrapingService := scrapingservices.NewJSONScrapingService(ctx, scrapingMapper)
	jsonScrapingService.Start()

	sjh := scrapingservices.NewSessionJobHandler(ctx, cfg,
		scrapingservices.WithMarketService("autovit", jsonScrapingService),
		scrapingservices.WithMarketService("mobile.de", collyScrapingService),
		scrapingservices.WithMarketService("bmw.de", collyScrapingService),
		scrapingservices.WithMarketService("autoscout", rodScrapingService),
		scrapingservices.WithMarketService("autotracknl", rodScrapingService),
		scrapingservices.WithMarketService("olx", jsonScrapingService),
		scrapingservices.WithMarketService("oferte_bmw", jsonScrapingService),
		scrapingservices.WithMarketService("tiriacauto", collyScrapingService),
		scrapingservices.WithMarketService("mercedes-benz.ro", jsonScrapingService),
		scrapingservices.WithMarketService("autoklass.ro", jsonScrapingService),
		scrapingservices.WithMarketService("mercedes-benz.de", jsonScrapingService),
	)
	sjh.Start()

	//https: //www.autotrack.nl/aanbod?minimumbouwjaar=2019&maximumkilometerstand=125000&brandstofsoorten=BENZINE&merkIds=7ccf5430-eafb-4042-82c0-43ce39ba1b02&modelIds=85e7360a-cee0-4ae0-85e0-0b595df99471&beschikbaarheidsStatus=beschikbaar&paginanummer=1&paginagrootte=30&sortering=PRIJS_OPLOPEND
	//https://www.autotrack.nl/aanbod?minimumbouwjaar=2019&maximumkilometerstand=125000&brandstofsoorten=BENZINE&modelIds=85e7360a-cee0-4ae0-85e0-0b595df99471&merkIds=7ccf5430-eafb-4042-82c0-43ce39ba1b02&                                    paginanummer=6&paginagrootte=30&sortering=PRIJS_OPLOPEND	sjh.StartWithoutMQ()

	//sjh := scrapingservices.NewSessionJobHandler(ctx, cfg)

	adapter := NewCriteriasJobsAdapter(cfg)
	go func() {
		//9,2023-11-20 00:06:39.350,2024-03-12 19:57:43.954,,autovit,www.autovit.ro,1
		//10,2023-11-20 00:06:39.350,2024-03-12 19:57:43.963,,webcar,www.webcar.ro,0
		//11,2023-11-20 00:06:39.350,2024-03-12 19:57:43.972,,mobile.de,www.mobile.de,1
		//12,2023-11-20 00:06:39.350,2024-03-12 19:57:43.978,,autoscout,www.autoscout24.ro,1
		//13,2023-11-20 00:06:39.350,2024-03-12 19:57:43.987,,autotracknl,www.autotrack.nl,0
		//14,2023-11-20 00:06:39.350,2024-03-12 19:57:43.996,,olx,www.olx.ro,1

		//markets := []uint{9, 11, 12, 13, 14, 15, 16, 17, 18, 19}
		//criterias := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
		markets := []uint{9}
		criterias := []uint{7}

		//markets := []uint{16}
		//criterias := []uint{2, 4, 5, 6, 12, 13}
		//criterias := []uint{22}
		//marketID := uint(9)

		allowedMarketAutoklassCriterias := []uint{8, 9, 24, 6, 13, 4, 1, 5, 27, 25, 28, 3, 10, 11, 19, 14}
		allowedMercedesBenzCriterias := []uint{3, 10, 11, 14, 19}

		allowedBMWDECriterias := []uint{1, 2, 5, 13}
		//allowedBMWCriterias := []uint{1, 2, 4, 5, 6, 12, 13}

		newSessionUUID := uuid.New()
		logSession, err := loggingService.CreateSession(newSessionUUID)
		if err != nil {
			panic(err)
		}

		for _, marketID := range markets {
			if marketID == 10 {
				marketID++
				continue
			}
			for _, criteriaID := range criterias {
				// criteria 7 volvo s90
				// criteria 6 bmw 7 series
				job := adapter.CreateJob(newSessionUUID, criteriaID, marketID)

				// do not scrape other brands for ofertebmw
				if job.Criteria.Brand != "bmw" && marketID == 15 {
					continue
				}

				//if marketID == 15 && !=inArrayUINT(criteriaID, allowedBMWCriterias) {
				//	continue
				//}

				if marketID == 20 && !inArrayUINT(criteriaID, allowedBMWDECriterias) {
					continue
				}

				if marketID == 18 && !inArrayUINT(criteriaID, allowedMarketAutoklassCriterias) {
					continue
				}

				if marketID == 17 || marketID == 19 {
					if !inArrayUINT(criteriaID, allowedMercedesBenzCriterias) {
						continue
					}
				}

				if job.Criteria.Brand != "mercedes-benz" && marketID == 17 {
					continue
				}

				//sjh.AddJobToMQ(*job)
				_, err := loggingService.CreateCriteriaLog(*logSession, *job)
				if err != nil {
					panic(err)
				}

				sjh.AddScrapingJob(*job)
			}
		}
		//done <- true
	}()

	<-done
	fmt.Println("exiting")

}

func inArrayUINT(str uint, list []uint) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func inArray(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
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
		AllowIncrementPage: true,
		SessionID:          sessionID,
		JobID:              uuid.New(),
		CriteriaID:         criteria.ID,
		MarketID:           marketID,
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
