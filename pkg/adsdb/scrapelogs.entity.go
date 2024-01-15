package adsdb

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScrapeLog struct {
	*gorm.Model
	SessionID   uuid.UUID
	JobID       uuid.UUID
	VisitURL    string
	Brand       string
	CarModel    string
	MarketName  string
	NumberOfAds int
	PageNumber  int
	IsLastPage  bool
	Error       string
}

type CriteriaLog struct {
	*gorm.Model
	SessionID   uuid.UUID
	Brand       string
	CarModel    string
	MarketName  string
	NumberOfAds int
	Error       string
	Success     bool
}

type SessionLog struct {
	*gorm.Model
	SessionID uuid.UUID
}
