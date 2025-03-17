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

type MarketsAndCriterias struct {
	Markets   []uint `json:"markets"`
	Criterias []uint `json:"criterias"`
}

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
	r.HandleFunc("/startMarketsAndCriterias", scrapeMarketCriteria(sessionService)).Methods("POST")

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
		res := Response{Data: "started scraping market and criteria"}
		resb, err := json.Marshal(&res)
		if err != nil {
			panic(err)
		}

		w.Write(resb)
	}
}

func scrapeMarketsAndCriterias(s *sessionstarter.SessionStarterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mac, err := getMarketsAndCriterias(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		s.ScrapeMarketsCriterias(mac.Markets, mac.Criterias)
		w.Header().Add("Access-Control-Allow-Origin", "*")
		type Response struct {
			Data string
		}
		res := Response{Data: "started scraping market and criteria"}
		resb, err := json.Marshal(&res)
		if err != nil {
			panic(err)
		}

		w.Write(resb)

	}
}

func getMarketsAndCriterias(r *http.Request) (MarketsAndCriterias, error) {
	var mac MarketsAndCriterias
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mac)
	if err != nil {
		return mac, err
	}
	return mac, nil
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
