package results

import (
	"carscraper/pkg/jobs"
)

type ResultsManager struct {
	Sessions []ProcessedSession
}

func (srm *ResultsManager) getSession(sessionID string) *ProcessedSession {

	for _, session := range srm.Sessions {
		if session.SessionID == sessionID {
			return &session
		}
	}
	ps := NewProcessedSession(sessionID)
	srm.Sessions = append(srm.Sessions, *ps)
	return ps
}

func NewResultsManager() *ResultsManager {
	return &ResultsManager{
		Sessions: make([]ProcessedSession, 0),
	}
}

func (srm *ResultsManager) AddPageResults(sessionID string, criteriaID uint, marketID uint, result jobs.AdsPageJobResult) {
	pageResult := NewPageResult(result.RequestedScrapingJob.Market.PageNumber, result.IsLastPage, result.Data)
	session := srm.getSession(sessionID)
	criteria := session.getCriteria(criteriaID)
	market := criteria.getMarket(marketID)
	market.AddPageResult(*pageResult)
}

func (srm *ResultsManager) GetAds(sessionID string, criteriaID uint, marketID uint) []jobs.Ad {
	mkt := srm.getSession(sessionID).getCriteria(criteriaID).getMarket(marketID)
	ads := make([]jobs.Ad, 0)
	for _, results := range *mkt.Results {
		pageAds := results.results
		ads = append(ads, *pageAds...)
	}
	return ads
}

func (srm *ResultsManager) isCompleteSession(sessionID string) bool {
	return srm.getSession(sessionID).isComplete()
}

func (srm *ResultsManager) isCompleteCriteria(sessionID string, criteriaID uint) bool {
	return srm.getSession(sessionID).getCriteria(criteriaID).isComplete()
}

func (srm *ResultsManager) isCompleteMarket(sessionID string, criteriaID uint, marketID uint) bool {
	return srm.getSession(sessionID).getCriteria(criteriaID).getMarket(marketID).isComplete()
}
