package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/scraping/scraper"
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
	)
	sessionService.Start()

	sjc := scraper.NewPageScrapingService(cfg, scraper.WithSimpleMessageQueueRepository(cfg))
	sjc.Start()
}
