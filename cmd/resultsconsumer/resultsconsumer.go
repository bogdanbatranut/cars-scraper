package main

import (
	"carscraper/pkg/adapters"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/notifications"
	"carscraper/pkg/repos"
	"carscraper/pkg/results"
)

func main() {
	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)

	resultsWriter := results.NewResultsWriter(adapters.NewAdsResultsAdapter(cfg), *repos.NewAdsRepository(cfg))
	rc := results.NewResultsReaderService(
		results.WithResultsMQRepository(cfg),
		results.WithLogger(cfg),
		results.WithTopicName(cfg),
		results.WithResultsWriter(*resultsWriter),
		results.WithRepo(repos.NewAdsRepository(cfg)),
		results.WithNotificationService(notifications.NewNotificationsService(cfg)),
	)
	rc.Start()
}
