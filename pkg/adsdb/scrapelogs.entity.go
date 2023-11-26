package adsdb

import (
	"gorm.io/gorm"
)

type ScrapeLog struct {
	*gorm.Model
	URL        *string
	CriteriaID uint
	MarketID   uint
	Success    bool
	Criteria   Criteria
	Market     Market
}
