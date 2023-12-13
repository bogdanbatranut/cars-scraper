package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/results"
)

func main() {
	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)

	resultsWriter := results.NewResultsWriter(results.NewAdsResultsAdapter(cfg), *results.NewResultsRepository(cfg))
	rc := results.NewResultsReaderService(results.WithResultsMQRepository(cfg), results.WithTopicName(cfg), results.WithResultsWriter(*resultsWriter))
	rc.Start()
}
