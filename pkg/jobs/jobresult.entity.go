package jobs

import (
	"carscraper/pkg/scraping/strategies"
)

type AdsPageJobResult struct {
	RequestedScrapingJob SessionJob
	PageNumber           int
	IsLastPage           bool
	Success              bool
	Data                 *[]strategies.Ad
}
