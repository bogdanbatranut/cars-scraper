package autovit

import (
	"carscraper/pkg/jobs"
	"carscraper/pkg/scraping/icollector"
	"encoding/json"
	"log"
)

type AutovitJSONAdapter struct {
}

func NewAutovitJSONAdapter() *AutovitJSONAdapter { return &AutovitJSONAdapter{} }

func (a AutovitJSONAdapter) GetAds(job jobs.SessionJob) icollector.AdsResults {
	var ads []jobs.Ad

	//fileNumberStr := strconv.Itoa(job.Market.PageNumber)
	autovitResults, getResultsERR := a.getJobResults(job)
	if getResultsERR != nil {
		panic(getResultsERR)
	}
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
	//a.loggingService.AddPageScrapeEntry(job, totalCount, job.Market.PageNumber, isLastPage, pageURL, getResultsERR)
	if getResultsERR != nil {
		panic(getResultsERR)
	}
	//return ads, isLastPage, nil
	return icollector.AdsResults{
		Ads:        &ads,
		IsLastPage: isLastPage,
		Error:      nil,
	}
}

func (s AutovitJSONAdapter) getJobResults(job jobs.SessionJob) (*AutovitGraphQLResponse, error) {
	r := NewRequest(job.Criteria)
	url := r.urlBuilder.GetPageURL(job.Market.PageNumber)
	responseBytes, err := r.DoRequest(url)
	if err != nil {
		return nil, err
	}
	var obj AutovitGraphQLResponse
	err = json.Unmarshal(responseBytes, &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}
