package marketadapters

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
)

type IJSONMarketAdsAdapter interface {
	GetAds(job jobs.SessionJob) icollector.AdsResults
}
