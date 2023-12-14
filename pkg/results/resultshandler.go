package results

import (
	"carscraper/pkg/jobs"
	"log"

	"github.com/google/uuid"
)

type SessionCriteriaMarketResultsHandler struct {
	results map[string]map[uint]map[uint]*MarketScrapingResults
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
				if results.adsInPage == nil {
					log.Printf("No results !!!")
					break
				}
				for _, resultsPages := range results.adsInPage {
					for _, res := range *resultsPages.results {
						log.Printf(" Make: %s Model: %s", res.Brand, res.Model)
					}
				}
			}
		}
	}
	log.Println("________________________________________________________________________________________________")
}

type MarketResults struct {
	MarketID     uint
	ResultsPages MarketScrapingResults
}

type MarketScrapingResults struct {
	adsInPage      []AdsInPage
	lastPageNumber *int
}

func (rp MarketScrapingResults) getNumOfExistingAds() int {
	exAds := rp.adsInPage
	return len(exAds)
}

func (rp MarketScrapingResults) getAds() []jobs.Ad {
	ads := []jobs.Ad{}
	if rp.adsInPage[0].results == nil {
		return nil
	}
	for _, adsInPage := range rp.adsInPage {
		ads = append(ads, *adsInPage.results...)
	}
	return ads
}

func NewResultsPages() *MarketScrapingResults {
	var adsArray = []AdsInPage{}
	lastPage := 0
	return &MarketScrapingResults{
		adsInPage:      adsArray,
		lastPageNumber: &lastPage,
	}

}

type AdsInPage struct {
	pageNumber int
	isLastPage bool
	results    *[]jobs.Ad
}

func AddPageResults(pageNumber int, isLastPage bool, ads *jobs.AdsPageJobResult, rep *MarketScrapingResults) {
	if rep.lastPageNumber == nil || *rep.lastPageNumber < 1 {
		if isLastPage {
			rep.lastPageNumber = &pageNumber
		}
	}
	aip := AdsInPage{
		pageNumber: pageNumber,
		results:    ads.Data,
	}
	rep.adsInPage = append(rep.adsInPage, aip)
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
	for _, adsInPage := range rp.adsInPage {
		if adsInPage.pageNumber == page {
			return true
		}
	}
	return false
}
