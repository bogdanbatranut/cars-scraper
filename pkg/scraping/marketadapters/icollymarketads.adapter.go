package marketadapters

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
)

type ICollyMarketAdsAdapter interface {
	GetAds(job jobs.SessionJob) icollector.AdsResults
}
