package markets

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
)

type IScrapingJob interface {
	Execute(job jobs.SessionJob) icollector.AdsResults
}
