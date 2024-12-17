package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/repos"
	"log"
	"time"
)

func main() {
	log.Println("starting BACKEND service...")

	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)

	//r := mux.NewRouter().StrictSlash(true)

	//criteriaRepo := repos.NewSQLCriteriaRepository(cfg)
	//marketsRepo := repos.NewSQLMarketsRepository(cfg)
	adsRepo := repos.NewAdsRepository(cfg)
	getCheapestToday(adsRepo)

	//r.HandleFunc("/test", test()).Methods("POST")

	//r.HandleFunc("/marketsAndCriterias", marketsAndCriterias(criteriaRepo)).Methods("POST")
	//
	//httpPort := cfg.GetString(amconfig.BackendServiceHTTPPort)
	//log.Printf("HTTP listening on port %s\n", httpPort)
	//err = http.ListenAndServe(fmt.Sprintf(":%s", httpPort), r)
	//errorshandler.HandleErr(err)

}

func getCheapestToday(adsRepo *repos.AdsRepository) {
	today := time.Now()

	ads, err := adsRepo.GetAll()
	if err != nil {
		panic(err)
	}

	for _, ad := range *ads {
		if isSameDay(ad.CreatedAt, today) {
			log.Println(ad.ID)
		}
	}
}

func isSameDay(t2 time.Time, t1 time.Time) bool {
	sameYear := t1.Year() == t2.Year()
	sameMonth := t1.Month() == t2.Month()
	sameDay := t1.Day() == t2.Day()

	return sameYear && sameMonth && sameDay

}
