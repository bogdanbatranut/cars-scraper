package main

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"carscraper/pkg/repos"
)

func main() {
	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)
	marketID := uint(9)
	muuid := "m_uuid"

	seller := adsdb.Seller{
		NameInMarket: "name in market",
		URLInMarket:  "url in market",
		OwnURL:       "ownurl",
	}

	price := adsdb.Price{
		Price:    200000,
		MarketID: marketID,
	}

	var prices []adsdb.Price
	prices = append(prices, price)

	ad := adsdb.Ad{
		Brand:      "test",
		CarModel:   "test",
		Year:       2019,
		Km:         1000,
		Fuel:       "diesel",
		MarketUUID: &muuid,
		Active:     true,
		Ad_url:     "www.adurl.com",
		MarketID:   9,
		//Market:     adsdb.Market{},
		Seller: seller,
		Prices: prices,
	}

	ads := []adsdb.Ad{}
	ads = append(ads, ad)

	seller2 := adsdb.Seller{
		NameInMarket: "name in market",
		URLInMarket:  "url in market",
		OwnURL:       "ownurl",
	}

	price2 := adsdb.Price{
		Price:    200000,
		MarketID: marketID,
	}
	muuid2 := "m_uuid2"

	var prices2 []adsdb.Price
	prices2 = append(prices2, price2)

	ad2 := adsdb.Ad{
		Brand:      "2test",
		CarModel:   "2test",
		Year:       2019,
		Km:         100000,
		Fuel:       "diesel",
		MarketUUID: &muuid2,
		Active:     true,
		Ad_url:     "www.adurl.com",
		MarketID:   10,
		//Market:     adsdb.Market{},
		Seller: seller2,
		Prices: prices2,
	}
	ads = append(ads, ad2)

	repo := repos.NewAdsRepository(cfg)
	repo.Upsert(ads)

}
