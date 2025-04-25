package scrapingservices

import (
	"carscraper/pkg/jobs"
)

type IScrapingService interface {
	AddJob(job jobs.SessionJob)
	GetResultsChannel() *chan jobs.AdsPageJobResult
	GetCurrentJobExecutionAvailabilityChannel() chan bool
}
