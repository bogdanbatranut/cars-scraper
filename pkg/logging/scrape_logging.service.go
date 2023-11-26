package logging

import (
	"carscraper/pkg/adsdb"

	"gorm.io/gorm"
)

type ScrapeLoggingService struct {
	db *gorm.DB
}

func NewScrapeLoggingService(db *gorm.DB) ScrapeLoggingService {
	return ScrapeLoggingService{
		db: db,
	}
}

func (sls ScrapeLoggingService) AddEntry(market adsdb.Market, criteria adsdb.Criteria, url string) *adsdb.ScrapeLog {
	logEntry := adsdb.ScrapeLog{
		URL:        &url,
		CriteriaID: criteria.ID,
		MarketID:   market.ID,
		Success:    false,
	}

	sls.db.Save(&logEntry)
	return &logEntry
}

func (sls ScrapeLoggingService) DoneWithSuccess(sl *adsdb.ScrapeLog) {
	sl.Success = true
	sls.db.Save(&sl)
}
