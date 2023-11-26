package jobs

import (
	"carscraper/pkg/scraping/strategies"
	"carscraper/pkg/scraping/urlbuilder"
)

type AdsPageJobResult struct {
	RequestedScrapingJob SessionJob
	PageURL              urlbuilder.PageURL
	PageNumber           int
	IsLastPage           bool
	Success              bool
	Data                 *[]strategies.Ad
}
