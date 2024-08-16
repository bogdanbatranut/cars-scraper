package adsdb

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionLog struct {
	*gorm.Model
	SessionID    uuid.UUID
	PageLogs     []PageLog
	CriteriaLogs []CriteriaLog
}

type CriteriaLog struct {
	*gorm.Model
	SessionID    uuid.UUID
	SessionLogID uint
	SessionLog   SessionLog
	CriteriaID   uint
	MarketID     uint
	//Market       Market
	//Criteria     Criteria
	Fuel        string
	Brand       string
	CarModel    string
	MarketName  string
	NumberOfAds int
	Success     bool
	Finished    bool
	PageLogs    []PageLog
}

type PageLog struct {
	*gorm.Model
	SessionLogID  uint
	SessionID     uuid.UUID
	SessionLog    SessionLog
	JobID         uuid.UUID
	CriteriaLogID uint
	CriteriaLog   CriteriaLog
	VisitURL      string
	Brand         string
	CarModel      string
	MarketName    string
	MarketID      uint
	//Market        Market
	NumberOfAds int
	PageNumber  int
	IsLastPage  bool
	Error       string
	Scraped     bool
	Consumed    bool
}
