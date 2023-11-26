package main

import (
	"carscraper/pkg/config"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/scraping"
	"log"
)

func main() {
	log.Println("starting sessionstarter service...")

	cfg, err := config.NewViperConfig()
	errorshandler.HandleErr(err)

	sessionService := scraping.NewSessionStarterService(
		scraping.WithSimpleMessageQueueRepository(cfg),
		scraping.WithCriteriaSQLRepository(cfg),
	)
	sessionService.Start()
}
