package results

import (
	"carscraper/pkg/adsdb"
)

type IResultsRepository interface {
	WriteResults(ads []adsdb.Ad) error
	GetAllAdsIDs(marketID uint, criteriaID uint) *[]uint
	DeleteAd(adID uint)
}
