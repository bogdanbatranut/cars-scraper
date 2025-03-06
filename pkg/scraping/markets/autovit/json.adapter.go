package autovit

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/jobs"
	"carscraper/pkg/logging"
	"carscraper/pkg/scraping/icollector"
	"encoding/json"
	"log"
)

type AutovitJSONAdapter struct {
	loggingService *logging.ScrapeLoggingService
}

func NewAutovitJSONAdapter(loggingService *logging.ScrapeLoggingService) *AutovitJSONAdapter {
	return &AutovitJSONAdapter{loggingService: loggingService}
}

func (a AutovitJSONAdapter) GetAds(job jobs.SessionJob) icollector.AdsResults {
	var ads []jobs.Ad

	criteriaLog, err := a.loggingService.GetCriteriaLog(job.SessionID, job.CriteriaID, job.MarketID)
	if err != nil {
		panic(err)
	}
	pageLog, err := a.loggingService.CreatePageLog(criteriaLog, job, "", job.Market.PageNumber)
	if err != nil {
		panic(err)
	}

	//fileNumberStr := strconv.Itoa(job.Market.PageNumber)
	autovitResults, getResultsERR := a.getJobResults(job, *pageLog)
	if getResultsERR != nil {
		a.loggingService.PageLogSetError(pageLog, getResultsERR.Error())
		return icollector.AdsResults{
			Ads:        nil,
			IsLastPage: true,
			Error:      getResultsERR,
		}
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

	err2 := a.loggingService.PageLogSetPageScraped(pageLog, len(ads), isLastPage)
	if err2 != nil {
		log.Println(err2.Error())
	}

	return icollector.AdsResults{
		Ads:        &ads,
		IsLastPage: isLastPage,
		Error:      nil,
	}
}

func (a AutovitJSONAdapter) getJobResults(job jobs.SessionJob, pageLog adsdb.PageLog) (*AutovitGraphQLResponse, error) {
	r := NewRequest(job.Criteria)
	url := r.urlBuilder.GetPageURL(job.Market.PageNumber)
	err := a.loggingService.PageLogSetVisitURL(&pageLog, url)
	if err != nil {
		log.Println(err.Error())
	}

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
