package main

import "carscraper/pkg/scraping"

func main() {

	sjc := scraping.NewPageScrapingService(scraping.WithSimpleMessageQueueRepository1())
	sjc.Start()
}
