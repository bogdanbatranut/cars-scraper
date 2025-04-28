package events

import "carscraper/pkg/adsdb"

type CreateEvent struct {
	Ad adsdb.Ad
}

type UpdateEvent struct {
	Ad adsdb.Ad
}

type UpdatePriceEvent struct {
	Ad adsdb.Ad
}

type DeleteEvent struct {
	Ad adsdb.Ad
}

type MinPriceCreatedEvent struct {
	Ad adsdb.Ad
}

type MinPriceUpdatedEvent struct {
	Ad adsdb.Ad
}
