package main

import "carscraper/pkg/notifications"

func main() {
	notificationsService := notifications.NewNotificationsService(nil)
	notificationsService.SendOpenAdNotification(600331, "Test notification")
}
