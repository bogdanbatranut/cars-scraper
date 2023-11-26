package adsdb

import "gorm.io/gorm"

type CriteriaMarkets struct {
	gorm.Model
	CriteriaID  uint
	MarketID    uint
	AllowScrape bool
}
