package results

import (
	"carscraper/pkg/jobs"
	"log"

	"github.com/google/uuid"
)

type SessionCriteriaMarketResultsHandler struct {
	results map[string]map[uint]map[uint]*MarketScrapingResults
}

type MarketScrapingResults struct {
	pageResults    []PageResult
	lastPageNumber *int
}

func NewSessionCriteriaMarketResults() *SessionCriteriaMarketResultsHandler {
	m := make(map[string]map[uint]map[uint]*MarketScrapingResults)
	return &SessionCriteriaMarketResultsHandler{
		results: m,
	}
}

func (scmr SessionCriteriaMarketResultsHandler) createMarketsMap(marketID uint, result jobs.AdsPageJobResult) map[uint]*MarketScrapingResults {
	rp := NewResultsPages()
	//rp.Add(result.PageNumber, result.IsLastPage, &result)

	var marketMap = map[uint]*MarketScrapingResults{}
	marketMap[marketID] = rp
	return marketMap
}

func (scmr SessionCriteriaMarketResultsHandler) createCriteriasMap(criteriaID uint, marketsMap map[uint]*MarketScrapingResults) map[uint]map[uint]*MarketScrapingResults {
	var criteriasMap = map[uint]map[uint]*MarketScrapingResults{}
	criteriasMap[criteriaID] = marketsMap
	return criteriasMap
}

func (scmr SessionCriteriaMarketResultsHandler) Add(sessionID uuid.UUID, criteriaID uint, marketID uint, result jobs.AdsPageJobResult) {
	sessionIDStr := sessionID.String()
	//scmr.results[sessionIDStr][criteriaID][marketID] = NewResultsPages()
	if scmr.results[sessionIDStr] == nil {
		marketsMap := scmr.createMarketsMap(marketID, result)
		criteriasMap := scmr.createCriteriasMap(criteriaID, marketsMap)
		scmr.results[sessionIDStr] = criteriasMap
	}
	if scmr.results[sessionIDStr][criteriaID] == nil {
		marketsMap := scmr.createMarketsMap(marketID, result)
		scmr.results[sessionIDStr][criteriaID] = marketsMap
	} else {
		if scmr.results[sessionIDStr][criteriaID][marketID] == nil {
			scmr.results[sessionIDStr][criteriaID][marketID] = NewResultsPages()
		}

	}
	if result.IsLastPage {
		log.Printf("Got last page for : make: %s model: %s", result.RequestedScrapingJob.Criteria.Brand, result.RequestedScrapingJob.Criteria.CarModel)
	}
	AddPageResults(result.RequestedScrapingJob.Market.PageNumber, result.IsLastPage, &result, scmr.results[sessionIDStr][criteriaID][marketID])

}

func (scmr SessionCriteriaMarketResultsHandler) Print() {
	log.Println("________________________________________________________________________________________________")
	for sessionID, criterias := range scmr.results {
		for criteriaID, markets := range criterias {
			for marketID, results := range markets {
				log.Printf("Session: %+v", sessionID)
				log.Printf("Criteria: %d", criteriaID)
				log.Printf("Market: %d", marketID)
				if results.pageResults == nil {
					log.Printf("No results !!!")
					break
				}
				for _, resultsPages := range results.pageResults {
					for _, res := range *resultsPages.results {
						log.Printf(" Make: %s Model: %s", res.Brand, res.Model)
					}
				}
			}
		}
	}
	log.Println("________________________________________________________________________________________________")
}

func (rp MarketScrapingResults) getAds() []jobs.Ad {
	ads := []jobs.Ad{}
	if rp.pageResults[0].results == nil {
		return nil
	}
	for _, adsInPage := range rp.pageResults {
		if *adsInPage.results == nil {
			log.Println("We have null pageResults.results")
			continue
		}
		ads = append(ads, *adsInPage.results...)
	}
	return ads
}

func NewResultsPages() *MarketScrapingResults {
	var adsArray = []PageResult{}
	lastPage := 0
	return &MarketScrapingResults{
		pageResults:    adsArray,
		lastPageNumber: &lastPage,
	}

}

type PageResult struct {
	pageNumber int
	isLastPage bool
	results    *[]jobs.Ad
}

func NewPageResult(pageNumber int, isLastPage bool, ads *[]jobs.Ad) *PageResult {
	return &PageResult{
		pageNumber: pageNumber,
		isLastPage: isLastPage,
		results:    ads,
	}
}

func AddPageResults(pageNumber int, isLastPage bool, ads *jobs.AdsPageJobResult, rep *MarketScrapingResults) {
	if rep.lastPageNumber == nil || *rep.lastPageNumber < 1 {
		if isLastPage {
			rep.lastPageNumber = &pageNumber
			log.Println("Got last page for a criteria...")
		}
	}
	aip := PageResult{
		pageNumber: pageNumber,
		results:    ads.Data,
	}
	rep.pageResults = append(rep.pageResults, aip)
}

func (rep MarketScrapingResults) IsComplete() bool {
	if rep.lastPageNumber == nil || *rep.lastPageNumber < 1 {
		return false
	}
	isComplete := true
	for i := 1; i <= *rep.lastPageNumber; i++ {
		foundPage := rep.findResultsPage(i)
		if !foundPage {
			return false
		}
	}

	return isComplete
}

func (rp MarketScrapingResults) findResultsPage(page int) bool {
	for _, adsInPage := range rp.pageResults {
		if adsInPage.pageNumber == page {
			return true
		}
	}
	return false
}
