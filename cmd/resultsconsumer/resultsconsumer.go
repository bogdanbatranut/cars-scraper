package main

import (
	"carscraper/pkg/adapters"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/events"
	"carscraper/pkg/notifications"
	"carscraper/pkg/repos"
	"carscraper/pkg/results"
)

func main() {
	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)
	notificationsService := notifications.NewNotificationsService(cfg)
	eventsListener := events.NewEventsListener(notificationsService)
	adsRepository := repos.NewAdsRepository(cfg, eventsListener)
	resultsWriter := results.NewResultsWriter(adapters.NewAdsResultsAdapter(cfg), *adsRepository)

	rc := results.NewResultsReaderService(
		results.WithResultsMQRepository(cfg),
		results.WithLogger(cfg),
		results.WithTopicName(cfg),
		results.WithResultsWriter(*resultsWriter),
		results.WithRepo(adsRepository),
		results.WithNotificationService(notificationsService),
		results.WithEventsListener(eventsListener),
	)
	rc.Start()
}
