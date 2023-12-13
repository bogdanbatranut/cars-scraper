package adsdb

import "gorm.io/gorm"

type Price struct {
	*gorm.Model
	Price    int
	AdID     uint
	MarketID uint
	Ad       Ad
	Market   Market
}
