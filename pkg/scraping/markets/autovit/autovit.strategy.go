package autovit

import (
	"carscraper/pkg/jobs"
	"encoding/json"
	"log"
)

type AutovitStrategy struct {
}

func NewAutovitStrategy() AutovitStrategy {
	return AutovitStrategy{}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (as AutovitStrategy) Execute(job jobs.SessionJob) ([]jobs.Ad, bool, error) {
	var ads []jobs.Ad

	//fileNumberStr := strconv.Itoa(job.Market.PageNumber)
	autovitResults := as.getJobResults(job)

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
	//isLastPage = true
	return ads, isLastPage, nil
}

func (s AutovitStrategy) getJobResults(job jobs.SessionJob) AutovitGraphQLResponse {
	r := NewRequest(job.Criteria)
	byteResults := r.GetPage(job.Market.PageNumber)
	var obj AutovitGraphQLResponse
	err := json.Unmarshal(byteResults, &obj)
	if err != nil {
		panic(err)
	}
	return obj
}
