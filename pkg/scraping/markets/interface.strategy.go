package markets

import "carscraper/pkg/jobs"

type IScrapingStrategy interface {
	Execute(job jobs.SessionJob) ([]jobs.Ad, bool, error)
}
