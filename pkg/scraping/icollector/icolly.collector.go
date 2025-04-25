package icollector

import (
	"carscraper/pkg/jobs"

	"github.com/gocolly/colly"
)

type ICollyCollector interface {
	GetCollyCollector(job jobs.SessionJob) (*colly.Collector, error)
}
