package main

import (
	"carscraper/pkg/config"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/sessionstarter"
	"log"
)

func main() {
	log.Println("starting sessionstarter service...")

	cfg, err := config.NewViperConfig()
	errorshandler.HandleErr(err)

	sessionService := sessionstarter.NewSessionStarterService(
		sessionstarter.WithSimpleMessageQueueRepository(cfg),
		sessionstarter.WithCriteriaSQLRepository(cfg),
	)
	sessionService.Start()
}
