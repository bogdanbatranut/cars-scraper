package autoscoutcollycollector

import (
	"carscraper/pkg/jobs"

	"github.com/gocolly/colly"
)

type AutoScoutCollyCollector struct {
}

func (collector AutoScoutCollyCollector) GetCollyCollector(job jobs.SessionJob) *colly.Collector {
	return colly.NewCollector(nil)
}
