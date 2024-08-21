package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/repos"
	"encoding/json"
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
	//ctx, cancel := context.WithCancel(context.Background())

	smqHost := cfg.GetString(amconfig.SMQURL)
	smqPort := cfg.GetString(amconfig.SMQHTTPPort)
	log.Printf("SMQ HOST: %s:%s ", smqHost, smqPort)
	smqr := repos.NewSimpleMessageQueueRepository(fmt.Sprintf("http://%s:%s", smqHost, smqPort))
	log.Println("Message Queue URL : ", fmt.Sprintf("http://%s:%s", smqHost, smqPort))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		log.Println("canceling")
		//cancel()
		done <- true
	}()

	fmt.Println("awaiting signal")

	adapter := DNewCriteriasJobsAdapter(cfg)
	loggingService := logging.NewScrapeLoggingService(cfg)

	go func() {
		//9,2023-11-20 00:06:39.350,2024-03-12 19:57:43.954,,autovit,www.autovit.ro,1
		//10,2023-11-20 00:06:39.350,2024-03-12 19:57:43.963,,webcar,www.webcar.ro,0
		//11,2023-11-20 00:06:39.350,2024-03-12 19:57:43.972,,mobile.de,www.mobile.de,1
		//12,2023-11-20 00:06:39.350,2024-03-12 19:57:43.978,,autoscout,www.autoscout24.ro,1
		//13,2023-11-20 00:06:39.350,2024-03-12 19:57:43.987,,autotracknl,www.autotrack.nl,0
		//14,2023-11-20 00:06:39.350,2024-03-12 19:57:43.996,,olx,www.olx.ro,1

		markets := []uint{20}
		//criterias := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29}
		//markets := []uint{14}
		criterias := []uint{1, 2, 5, 13}

		//markets := []uint{16}
		//criterias := []uint{2, 4, 5, 6, 12, 13}
		//criterias := []uint{22}
		//marketID := uint(9)

		allowedMarketAutoklassCriterias := []uint{8, 9, 24, 6, 13, 4, 1, 5, 27, 25, 28, 3, 10, 11, 19, 14}
		allowedMercedesBenzCriterias := []uint{3, 10, 11, 14, 19}
		allowedBMWDECriterias := []uint{1, 2, 5, 6, 13}
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
				job := adapter.DCreateJob(newSessionUUID, criteriaID, marketID)

				// do not scrape other brands for ofertebmw
				if job.Criteria.Brand != "bmw" && marketID == 15 {
					continue
				}

				if marketID == 20 && !DinArrayUINT(criteriaID, allowedBMWDECriterias) {
					continue
				}

				//if marketID == 15 && !=inArrayUINT(criteriaID, allowedBMWCriterias) {
				//	continue
				//}

				if marketID == 18 && !DinArrayUINT(criteriaID, allowedMarketAutoklassCriterias) {
					continue
				}

				if marketID == 17 || marketID == 19 {
					if !DinArrayUINT(criteriaID, allowedMercedesBenzCriterias) {
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
				jobBytes, err := json.Marshal(&job)
				if err != nil {
					panic(err)
				}
				smqr.PutMessage("jobs", jobBytes)
			}
		}
		//done <- true
	}()

	<-done
	fmt.Println("exiting")

}

func DinArrayUINT(str uint, list []uint) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func DinArray(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

type DCriteriasJobsAdapter struct {
	criteriasRepo *repos.SQLCriteriaRepository
	marketsRepo   *repos.SQLMarketsRepository
}

func DNewCriteriasJobsAdapter(config amconfig.IConfig) *DCriteriasJobsAdapter {
	return &DCriteriasJobsAdapter{
		criteriasRepo: repos.NewSQLCriteriaRepository(config),
		marketsRepo:   repos.NewSQLMarketsRepository(config),
	}
}
func (adapter DCriteriasJobsAdapter) DCreateJob(sessionID uuid.UUID, criteriaID uint, marketID uint) *jobs.SessionJob {
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
