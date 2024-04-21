package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/sessionstarter"
	"log"
)

func main() {
	log.Println("starting  test sessionstarter service...")

	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)

	sessionService := sessionstarter.NewSessionStarterService(
		sessionstarter.WithSimpleMessageQueueRepository(cfg),
		sessionstarter.WithCriteriaSQLRepository(cfg),
		sessionstarter.WithLogging(cfg),
	)
	sessionService.Start()

	sjc := NewPageScrapingService(cfg, temp.WithSimpleMessageQueueRepository(cfg))
	sjc.Start()
}
