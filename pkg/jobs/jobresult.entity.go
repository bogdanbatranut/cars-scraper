package jobs

import "fmt"

type AdsPageJobResult struct {
	RequestedScrapingJob SessionJob
	PageNumber           int
	IsLastPage           bool
	Success              bool
	Data                 *[]Ad
}

func (res AdsPageJobResult) GetTopic() string {
	return fmt.Sprintf("results-%s-%d", res.RequestedScrapingJob.SessionID.String(), res.RequestedScrapingJob.CriteriaID)
}
