package repos

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IAdsRepository interface {
	GetAll() (*[]adsdb.Ad, error)
	GetAllAdsIDs(marketID uint, criteriaID uint) *[]uint
	GetAdsForCriteria(criteriaID uint, markets []string) *[]adsdb.Ad
	Upsert(ads []adsdb.Ad) (*[]uint, error)
	DeleteAd(adID uint)
	DeletePrice(priceID uint)
	GetAdPrices(adID uint) []adsdb.Price
}

type AdsRepository struct {
	db *gorm.DB
}

func NewAdsRepository(cfg amconfig.IConfig) *AdsRepository {
	databaseName := cfg.GetString(amconfig.AppDBName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &AdsRepository{
		db: db,
	}
}

func (r AdsRepository) GetAll() (*[]adsdb.Ad, error) {
	var ads []adsdb.Ad
	tx := r.db.Find(&ads)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &ads, nil
}

func (r AdsRepository) Upsert(ads []adsdb.Ad) (*[]uint, error) {
	adsIds := []uint{}
	r.db.Session(&gorm.Session{FullSaveAssociations: true})
	transactionErr := r.db.Transaction(func(tx *gorm.DB) error {
		var err error
		for _, foundAd := range ads {
			foundAdPrice := foundAd.Prices[0].Price
			foundMarketUUID := foundAd.MarketUUID
			foundAdKm := foundAd.Km

			tx = r.db.FirstOrCreate(&foundAd, adsdb.Ad{MarketUUID: foundMarketUUID}, adsdb.Ad{Prices: foundAd.Prices}, adsdb.Ad{SellerID: foundAd.SellerID})
			if tx.Error != nil {
				err = tx.Error
			}

			if foundAd.Km != foundAdKm {
				foundAd.Km = foundAdKm
				r.db.Save(&foundAd)
			}

			// get last price id db
			var lastExistingPrice adsdb.Price
			tx = r.db.Last(&lastExistingPrice, adsdb.Price{AdID: foundAd.ID})
			if tx.Error != nil {
				err = tx.Error
			}

			if lastExistingPrice.Price != foundAdPrice {
				// insert new price
				tx = r.db.Create(&adsdb.Price{Price: foundAdPrice, AdID: foundAd.ID, MarketID: foundAd.MarketID})
				if tx.Error != nil {
					err = tx.Error
				}
			}

			if foundAd.MarketID == 11 {
				if !strings.Contains(foundAd.Ad_url, "www.mobile.de") {
					tx := r.db.Model(&foundAd).Update("ad_url", fmt.Sprintf("https://www.mobile.de%s", foundAd.Ad_url))
					if tx.Error != nil {
						err = tx.Error
					}
				}

			}

			if foundAd.MarketID == 12 {
				if !strings.Contains(foundAd.Ad_url, "www.autoscout24.ro") {
					tx := r.db.Model(&foundAd).Update("ad_url", fmt.Sprintf("https://www.autoscout24.ro%s", foundAd.Ad_url))
					if tx.Error != nil {
						err = tx.Error
					}
				}
			}

			//price := foundAd.Prices[0]
			//price.AdID = foundAd.ID
			//tx = r.db.Debug().FirstOrCreate(&price, adsdb.Price{Price: foundAdPrice}, adsdb.Price{AdID: foundAd.ID})
			//if tx.Error != nil {
			//	err = tx.Error
			//}
			adsIds = append(adsIds, foundAd.ID)
		}
		if err != nil {
			return err
		}
		return nil
	})
	if transactionErr != nil {
		return nil, transactionErr
	}
	return &adsIds, transactionErr
}

func (r AdsRepository) DeleteAd(adID uint) {
	var ad adsdb.Ad
	r.db.Model(adsdb.Ad{}).Delete(&ad, adID)
}

func (r AdsRepository) DeletePrice(priceID uint) {
	model := gorm.Model{
		ID: priceID,
	}
	price := adsdb.Price{
		Model: &model,
	}
	r.db.Model(adsdb.Price{}).Unscoped().Delete(&price)
}

func (r AdsRepository) GetAllAdsIDs(marketID uint, criteriaID uint) *[]uint {
	var ads []adsdb.Ad
	var adsIDs []uint
	r.db.Model(adsdb.Ad{}).Select("id").Where("market_id = ? AND criteria_id = ?", marketID, criteriaID).Find(&ads)
	for _, ad := range ads {
		adsIDs = append(adsIDs, ad.ID)
	}
	return &adsIDs
}

func (r AdsRepository) GetAdsForCriteria(criteriaID uint, markets []string) *[]adsdb.Ad {
	var ads []adsdb.Ad
	r.db.Debug().Preload("Prices").Preload("Market").Where("criteria_id = ?", criteriaID).Where("market_id", markets).Find(&ads)
	return &ads
}

func (r AdsRepository) GetAdPrices(adID uint) []adsdb.Price {
	var prices []adsdb.Price
	price := adsdb.Price{AdID: adID}
	tx := r.db.Find(&prices, price)
	if tx.Error != nil {
		panic(tx.Error)
	}
	return prices
}
