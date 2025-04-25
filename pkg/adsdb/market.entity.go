package adsdb

import (
	"gorm.io/gorm"
)

type Market struct {
	gorm.Model
	Name         string
	URL          string
	AllowProcess bool
	Ads          []Ad
	Sellers      []Seller    `gorm:"many2many:seller_markets;"`
	Criterias    *[]Criteria `gorm:"many2many:criteria_markets;"`
}

func (m Market) TableName() string {
	return "automall.markets"
}
