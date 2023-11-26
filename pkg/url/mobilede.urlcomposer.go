package url

import "carscraper/pkg/adsdb"

type MobileDeURLComposer struct {
}

func NewMobileDeURLComposer() *MobileDeURLComposer {
	return &MobileDeURLComposer{}
}

func (ac MobileDeURLComposer) Create(criteria adsdb.Criteria) string {
	return "Implement this"
}
