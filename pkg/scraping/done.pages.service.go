package scraping

import (
	"carscraper/pkg/jobs"

	"github.com/google/uuid"
)

type SessionMarketCriteriaScrapingStatus struct {
	message string
}

func (dps DonePagesService) addJob(job jobs.PageToScrapeJob) SessionMarketCriteriaScrapingStatus {
	if job.PageURL.PageNumber < 0 {
		dps.pages[job.SessionID][job.MarketID][job.CriteriaID] = nil
		return SessionMarketCriteriaScrapingStatus{message: "success"}
	}
	existingJobs := append(*dps.pages[job.SessionID][job.MarketID][job.CriteriaID], job)
	for _, page := range existingJobs {
		if page.PageURL.PageNumber < 0 {
			dps.pages[job.SessionID][job.MarketID][job.CriteriaID] = nil
			return SessionMarketCriteriaScrapingStatus{message: "success"}
		}
	}
	return SessionMarketCriteriaScrapingStatus{message: "incomplete"}
}

func (dps DonePagesService) DeleteCriteria(sessionID uuid.UUID, marketID uint, criteria uint) {
	delete(dps.pages[sessionID][marketID], criteria)
}

func (dps DonePagesService) DeleteMarket(sessionID uuid.UUID, marketID uint) {
	delete(dps.pages[sessionID], marketID)
}

func (dps DonePagesService) DeleteSession(sessionID uuid.UUID) {
	delete(dps.pages, sessionID)
}

type DonePagesService struct {
	// service for keeping track of complete market criterias...
	// once we have valid completed market criteria, we can process those results...
	// sessionstarter, market, criteria -> jobs
	pages map[uuid.UUID]map[uint]map[uint]*[]jobs.PageToScrapeJob
}

func NewDonePagesService() *DonePagesService {
	return &DonePagesService{
		pages: make(map[uuid.UUID]map[uint]map[uint]*[]jobs.PageToScrapeJob),
	}
}

func (dps DonePagesService) PutJob(job jobs.PageToScrapeJob) string {
	status := dps.addJob(job)
	if status.message == "success" {
		dps.DeleteCriteria(job.SessionID, job.MarketID, job.CriteriaID)
		if len(dps.pages[job.SessionID][job.MarketID]) == 0 {
			dps.DeleteMarket(job.SessionID, job.MarketID)
			if len(dps.pages[job.SessionID]) == 0 {
				dps.DeleteSession(job.SessionID)
			}
		}
		return "done"
	}
	return "incomplete"
}
