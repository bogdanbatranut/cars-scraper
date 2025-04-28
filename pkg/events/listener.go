package events

import (
	"carscraper/pkg/notifications"
	"log"
)

type EventsListener struct {
	notificationService *notifications.NotificationsService
}

func NewEventsListener(service *notifications.NotificationsService) *EventsListener {
	return &EventsListener{
		notificationService: service,
	}
}

func (el *EventsListener) Fire(event interface{}) {
	switch e := event.(type) {
	case CreateEvent:
		el.handleCreate(e)
	case UpdateEvent:
		el.handleUpdate(e)
	case UpdatePriceEvent:
		el.handleUpdatePrice(e)
	case DeleteEvent:
		el.handleDelete(e)
	case MinPriceUpdatedEvent:
		el.handleMinPriceWhenUpdated(e)
	case MinPriceCreatedEvent:
		el.handleMinPriceCreated(e)
	default:
		panic("unknown event type")
	}

}

func (el *EventsListener) handleCreate(event CreateEvent) {
	// Handle create event

}
func (el *EventsListener) handleUpdate(event UpdateEvent) {

}
func (el *EventsListener) handleUpdatePrice(event UpdatePriceEvent) {
	// Handle update price event
	lastPrice := event.Ad.Prices[len(event.Ad.Prices)-1].Price
	secondToLastPrice := event.Ad.Prices[len(event.Ad.Prices)-2].Price
	if lastPrice < secondToLastPrice {
		err := el.notificationService.SendPriceDecreaseNotification(event.Ad)
		if err != nil {
			panic(err)
		}
	}

}
func (el *EventsListener) handleDelete(event DeleteEvent) {
	err := el.notificationService.SendDeleteAdNotification(event.Ad)
	if err != nil {
		log.Println(err)
		return
	}
}

func (el *EventsListener) handleMinPriceCreated(event MinPriceCreatedEvent) {
	err := el.notificationService.SendMinPriceCreatedNotification(event.Ad)
	if err != nil {
		log.Println(err)
		return
	}
}

func (el *EventsListener) handleMinPriceWhenUpdated(event MinPriceUpdatedEvent) {
	err := el.notificationService.SendNewMinPrice(event.Ad)
	if err != nil {
		log.Println(err)
		return
	}
}
