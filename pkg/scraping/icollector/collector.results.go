package icollector

import "carscraper/pkg/jobs"

type AdsResults struct {
	Ads        *[]jobs.Ad
	IsLastPage bool
	Error      error
}
