package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/sessionstarter"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	log.Println("starting sessionstarter service...")

	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)

	sessionService := sessionstarter.NewSessionStarterService(
		sessionstarter.WithSimpleMessageQueueRepository(cfg),
		sessionstarter.WithCriteriaSQLRepository(cfg),
		sessionstarter.WithMarketsSQLRepository(cfg),
		sessionstarter.WithLogging(cfg),
	)

	log.Println("sessionstarter service initialization...")
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/start", start(sessionService)).Methods("POST")
	r.HandleFunc("/startMarket/{marketID}", scrapeMarket(sessionService)).Methods("POST")
	r.HandleFunc("/startMarketCriteria/{marketID}/{criteriaID}", scrapeMarketCriteria(sessionService)).Methods("POST")

	appPort := cfg.GetString(amconfig.SessionStarterHTTPPort)
	log.Printf("HTTP listening on port %s\n", appPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", appPort), r)
	errorshandler.HandleErr(err)

}

func scrapeMarket(s *sessionstarter.SessionStarterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		marketId := vars["marketID"]

		s.ScrapeMarket(marketId)
		w.Header().Add("Access-Control-Allow-Origin", "*")
		type Response struct {
			Data string
		}
		res := Response{Data: "started scraping market"}
		resb, err := json.Marshal(&res)
		if err != nil {
			panic(err)
		}

		w.Write(resb)
	}
}

func scrapeMarketCriteria(s *sessionstarter.SessionStarterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		marketId := vars["marketID"]
		criteriaId := vars["criteriaID"]

		s.ScrapeMarketCriteria(marketId, criteriaId)
		w.Header().Add("Access-Control-Allow-Origin", "*")
		type Response struct {
			Data string
		}
		res := Response{Data: "started scraping market"}
		resb, err := json.Marshal(&res)
		if err != nil {
			panic(err)
		}

		w.Write(resb)
	}
}

func start(s *sessionstarter.SessionStarterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Start()
		w.Header().Add("Access-Control-Allow-Origin", "*")
		type Response struct {
			Data string
		}
		res := Response{Data: "started scraping session starter"}
		resb, err := json.Marshal(&res)
		if err != nil {
			panic(err)
		}

		w.Write(resb)
	}
}
