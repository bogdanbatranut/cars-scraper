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
	SessionID    uuid.UUID `gorm:"index"`
	SessionLogID uint
	SessionLog   SessionLog
	CriteriaID   uint `gorm:"index"`
	MarketID     uint `gorm:"index"`
	Fuel         string
	Brand        string
	CarModel     string
	MarketName   string
	NumberOfAds  int
	Success      bool
	Finished     bool
	PageLogs     []PageLog
}

type PageLog struct {
	*gorm.Model
	SessionLogID  uint      `gorm:"index"`
	SessionID     uuid.UUID `gorm:"index"`
	SessionLog    SessionLog
	JobID         uuid.UUID `gorm:"index"`
	CriteriaLogID uint      `gorm:"index"`
	CriteriaLog   CriteriaLog
	VisitURL      string
	Brand         string
	CarModel      string
	MarketName    string
	MarketID      uint
	//Market        Market
	NumberOfAds int
	PageNumber  int `gorm:"index"`
	IsLastPage  bool
	Error       string
	Scraped     bool
	Consumed    bool
}
