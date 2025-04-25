package main

import "carscraper/pkg/notifications"

func main() {
	notificationsService := notifications.NewNotificationsService(nil)
	notificationsService.SendMinPriceCreatedNotification(600331, "Test notification")
}
