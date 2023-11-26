package url

import "carscraper/pkg/adsdb"

type IURLComposer interface {
	Create(criteria adsdb.Criteria) string
}
