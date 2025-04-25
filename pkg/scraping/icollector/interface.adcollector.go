package icollector

import (
	"carscraper/pkg/jobs"
)

type IAdCollector interface {
	GetAds(url string, pageNumber int, criteria jobs.Criteria) AdsResults
}
