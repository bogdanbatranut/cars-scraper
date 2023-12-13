package adsdb

import "gorm.io/gorm"

type Seller struct {
	*gorm.Model
	NameInMarket string
	URLInMarket  string
	OwnURL       string
	Ads          []Ad
	Markets      []Market `gorm:"many2many:seller_markets;"`
}

type SellerMarkets struct {
	SellerID uint
	MarketID uint
}
