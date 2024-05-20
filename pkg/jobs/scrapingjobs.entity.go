package jobs

import (
	"fmt"

	"github.com/google/uuid"
)

type Criteria struct {
	Brand    string
	CarModel string
	YearFrom *int
	YearTo   *int
	Fuel     string
	KmFrom   *int
	KmTo     *int
}

type Market struct {
	Name       string
	PageNumber int
}

type Session struct {
	SessionID uuid.UUID
	Jobs      []SessionJob
}

type SessionJob struct {
	SessionID  uuid.UUID
	JobID      uuid.UUID
	CriteriaID uint
	MarketID   uint
	Criteria   Criteria
	Market     Market
}

func (job SessionJob) ToString() string {
	return fmt.Sprintf("JOB | Market: %s | Brand: %s | Model: %s | Fuel: %s | Page : %d",
		job.Market.Name, job.Criteria.Brand, job.Criteria.CarModel, job.Criteria.Fuel, job.Market.PageNumber)
}

type ScrapeResult struct {
	RequestedJob Session
	Results      *AdsPageJobResult
}

type PageToScrapeJob struct {
	ID         uuid.UUID
	SessionID  uuid.UUID
	MarketID   uint
	CriteriaID uint
	PageURL    PageURL
	Visited    bool
}
