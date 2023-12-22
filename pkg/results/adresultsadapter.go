package results

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/jobs"
	"carscraper/pkg/repos"
	"errors"

	"gorm.io/gorm"
)

type IAdsResultsAdapter interface {
	ToActiveDBAd(jobs.Ad, uint, uint) (*adsdb.Ad, error)
}

type AdsResultsAdapter struct {
	sellersRepo repos.ISellerRepository
}

func NewAdsResultsAdapter(cfg amconfig.IConfig) *AdsResultsAdapter {
	return &AdsResultsAdapter{sellersRepo: repos.NewSellerRepository(cfg)}
}

func (adapter AdsResultsAdapter) ToActiveDBAd(ad jobs.Ad, marketID uint, criteriaID uint) (*adsdb.Ad, error) {
	// find seller
	seller, err := adapter.sellersRepo.GetByURLInMarket(*ad.SellerMarketURL)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			seller = &adsdb.Seller{}
			seller.NameInMarket = *ad.SellerName
			seller.URLInMarket = *ad.SellerMarketURL
			s, err := adapter.sellersRepo.AddSeller(*seller)
			seller = s
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	price := adsdb.Price{
		Model:    &gorm.Model{},
		Price:    ad.Price,
		MarketID: marketID,
		Ad:       adsdb.Ad{},
		Market:   adsdb.Market{},
	}

	var prices []adsdb.Price
	prices = append(prices, price)
	return &adsdb.Ad{
		Brand:      ad.Brand,
		CarModel:   ad.Model,
		Year:       ad.Year,
		Km:         ad.Km,
		Fuel:       ad.Fuel,
		MarketUUID: &ad.AdID,
		Active:     true,
		Ad_url:     ad.Ad_url,
		MarketID:   marketID,
		CriteiaID:  criteriaID,
		SellerID:   seller.ID,
		// TODO Implement prices repo
		Prices: prices,
	}, nil
}
