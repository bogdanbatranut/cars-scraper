package age

import (
	"carscraper/pkg/adsdb"
	"time"
)

func CalculateAdAge(ad adsdb.Ad) int {
	currentTime := time.Now()
	adFirstSeenTime := ad.CreatedAt
	diff := currentTime.Sub(adFirstSeenTime)
	return int(diff.Hours() / 24)
}
