package adsdb

import (
	"gorm.io/gorm"
)

type Market struct {
	gorm.Model
	Name       string
	URL        string
	Ads        []Ad
	Sellers    []Seller    `gorm:"many2many:seller_markets;"`
	Criterias  *[]Criteria `gorm:"many2many:criteria_markets;"`
	ScrapeLogs *[]ScrapeLog
}
