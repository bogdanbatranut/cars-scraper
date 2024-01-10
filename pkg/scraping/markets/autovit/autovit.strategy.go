package autovit

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"encoding/json"
	"log"
)

type AutovitStrategy struct {
	loggingService logging.ScrapeLoggingService
}

func NewAutovitStrategy(logginService logging.ScrapeLoggingService) AutovitStrategy {
	return AutovitStrategy{
		loggingService: logginService,
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (as AutovitStrategy) Execute(job jobs.SessionJob) ([]jobs.Ad, bool, error) {
	var ads []jobs.Ad

	//fileNumberStr := strconv.Itoa(job.Market.PageNumber)
	autovitResults, pageURL, getResultsERR := as.getJobResults(job)

	for _, carData := range autovitResults.Data.AdvertSearch.Edges {
		ad := carData.Node.ToAd()

		ads = append(ads, *ad)

	}
	log.Printf("Got Autovit results : %d", len(ads))

	isLastPage := false
	totalCount := autovitResults.Data.AdvertSearch.TotalCount
	offSet := autovitResults.Data.AdvertSearch.PageInfo.CurrentOffset

	if totalCount-offSet <= autovitResults.Data.AdvertSearch.PageInfo.PageSize {
		isLastPage = true
	}
	//pageSize := autovitResults.Data.AdvertSearch.PageInfo.PageSize
	//pageNumber := offSet / pageSize
	as.loggingService.AddPageScrapeEntry(job, totalCount, job.Market.PageNumber, isLastPage, pageURL, getResultsERR)
	if getResultsERR != nil {
		panic(getResultsERR)
	}
	//isLastPage = true
	return ads, isLastPage, nil
}

func (s AutovitStrategy) getJobResults(job jobs.SessionJob) (*AutovitGraphQLResponse, string, error) {
	r := NewRequest(job.Criteria)
	byteResults, pageURL, err := r.GetPage(job.Market.PageNumber)
	if err != nil {
		return nil, pageURL, err
	}
	var obj AutovitGraphQLResponse
	err = json.Unmarshal(byteResults, &obj)
	if err != nil {
		return nil, pageURL, err
	}
	return &obj, pageURL, nil
}
