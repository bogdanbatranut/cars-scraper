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
		sessionstarter.WithLogging(cfg),
	)

	log.Println("sessionstarter service init...")
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/start", start(sessionService)).Methods("POST")

	appPort := cfg.GetString(amconfig.SessionStarterHTTPPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", appPort), r)
	errorshandler.HandleErr(err)
	log.Printf("HTTP listening on port %s\n", appPort)
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
