package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/scraping/scraper"
	"log"
)

func main() {

	log.Println("starting page scraping service...")

	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)

	sjc := scraper.NewPageScrapingService(cfg, scraper.WithSimpleMessageQueueRepository(cfg))
	sjc.Start()
}
