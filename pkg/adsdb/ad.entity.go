package adsdb

import "gorm.io/gorm"

type Ad struct {
	*gorm.Model
	Title        *string `gorm:"column:title"`
	Brand        string  `gorm:"not null"`
	CarModel     string  `gorm:"not null"`
	Year         int     `gorm:"not null"`
	Km           int     `gorm:"not null"`
	Fuel         string  `gorm:"not null"`
	MarketUUID   *string `gorm:"not null"`
	Active       bool    `gorm:"not null"`
	Ad_url       string  `gorm:"not null"`
	MarketID     uint    `gorm:"not null"`
	CriteriaID   uint    `gorm:"not null"`
	SellerID     uint    `gorm:"not null"`
	CurrentPrice *int    `gorm:"column:current_price"`
	Thumbnail    *string
	Market       Market
	Seller       Seller
	Prices       []Price
	Followed     bool `gorm:"not null"`
	//DealerMarketPrices []DealerMarketPrice `gorm:"many2many:dealer_market_prices;"`
}
