package main

import (
	"carscraper/pkg/config"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/scraping"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	log.Println("starting sessionstarter service...")

	cfg, err := config.NewViperConfig()
	errorshandler.HandleErr(err)

	sessionService := scraping.NewSessionStarterService(
		scraping.WithSimpleMessageQueueRepository(cfg),
		scraping.WithCriteriaSQLRepository(cfg),
	)

	log.Println("sessionstarter service init...")
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/start", start(sessionService)).Methods("POST")

	appPort := cfg.GetString(config.SessionStarterPort)
	err = http.ListenAndServe(fmt.Sprintf(":%s", appPort), r)
	errorshandler.HandleErr(err)
	log.Printf("HTTP listening on port %s\n", appPort)
}

func start(s *scraping.SessionStarterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Start()
		w.Write([]byte("started scraping session starter"))
	}
}
