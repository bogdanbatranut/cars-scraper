package adsdb

import "gorm.io/gorm"

type Ad struct {
	*gorm.Model
	Brand      string  `gorm:"not null"`
	CarModel   string  `gorm:"not null"`
	Year       int     `gorm:"not null"`
	Km         int     `gorm:"not null"`
	Fuel       string  `gorm:"not null"`
	MarketUUID *string `gorm:"not null"`
	Active     bool    `gorm:"not null"`
	Ad_url     string  `gorm:"not null"`
	MarketID   uint    `gorm:"not null"`
	CriteiaID  uint    `gorm:"not null"`
	SellerID   uint    `gorm:"not null"`
	Market     Market
	Seller     Seller
	Prices     []Price
	//DealerMarketPrices []DealerMarketPrice `gorm:"many2many:dealer_market_prices;"`
}
