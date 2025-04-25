package autotrackcollycollector

import (
	"carscraper/pkg/jobs"

	"github.com/gocolly/colly"
)

type AutotrackCollyCollector struct{}

func NewAutotrackCollyCollector() *AutotrackCollyCollector {
	return &AutotrackCollyCollector{}
}

func (collector AutotrackCollyCollector) GetCollyCollector(job jobs.SessionJob) *colly.Collector {
	return colly.NewCollector()
}
