package repos

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"fmt"
	"math"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int   `json:"limit,omitempty;query:limit"`
	Page       int   `json:"page,omitempty;query:page"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
	//Rows       interface{} `json:"rows"`
}

type IAdsRepository interface {
	GetAll() (*[]adsdb.Ad, error)
	GetAllAdsIDs(marketID uint, criteriaID uint) *[]uint
	GetAdsForCriteria(criteriaID uint, markets []string, minKm *int, maxKm *int, minPrice *int, maxPrice *int) *[]adsdb.Ad
	GetAdsForCriteriaPaginated(pagination *Pagination, criteriaID uint, markets []string, minKm *int, maxKm *int, minPrice *int, maxPrice *int) (*[]adsdb.Ad, *Pagination)
	Upsert(ads []adsdb.Ad) (*[]uint, error)
	DeleteAd(adID uint)
	DeletePrice(priceID uint)
	GetAdPrices(adID uint) []adsdb.Price
	UpdateCurrentPrice(adID uint)
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
			thumbnail := foundAd.Thumbnail
			foundAdPrice := foundAd.Prices[0].Price
			foundMarketUUID := foundAd.MarketUUID
			foundAdKm := foundAd.Km
			currentPrice := foundAd.CurrentPrice

			tx = r.db.FirstOrCreate(&foundAd, adsdb.Ad{MarketUUID: foundMarketUUID}, adsdb.Ad{Prices: foundAd.Prices}, adsdb.Ad{SellerID: foundAd.SellerID})
			if tx.Error != nil {
				err = tx.Error
			}

			if *foundAd.CurrentPrice != *currentPrice {
				foundAd.CurrentPrice = currentPrice
				r.db.Save(&foundAd)
			}

			//if foundAd.Thumbnail == nil && thumbnail != nil {
			foundAd.Thumbnail = thumbnail
			r.db.Save(&foundAd)
			//}
			//
			//if foundAd.Thumbnail != nil && strings.HasPrefix(*foundAd.Thumbnail, "data:image") && *foundAd.Thumbnail != *thumbnail {
			//	foundAd.Thumbnail = thumbnail
			//	r.db.Save(&foundAd)
			//}

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

func (r AdsRepository) GetAdsForCriteria(criteriaID uint, markets []string, minKm *int, maxKm *int, minPrice *int, maxPrice *int) *[]adsdb.Ad {
	var ads []adsdb.Ad
	r.db.Preload("Prices").Preload("Market").Where("criteria_id = ?", criteriaID).Where("market_id", markets).Where("current_price <= ?", maxPrice).Where("current_price >= ? ", minPrice).Find(&ads)
	return &ads
}

func withLimitAndOffset(pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}

func WithFilterConditions(db *gorm.DB, criteriaID uint, markets []string, minKm *int, maxKm *int, minPrice *int, maxPrice *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload("Prices").Preload("Market").Where("criteria_id = ?", criteriaID).Where("market_id", markets).Where("current_price <= ?", maxPrice).Where("current_price >= ? ", minPrice)
	}
}

//func WithFilterConditions(db *gorm.DB, criteriaID uint, markets []string, minKm *int, maxKm *int, minPrice *int, maxPrice *int) *gorm.DB {
//	return db.Preload("Prices").Preload("Market").Where("criteria_id = ?", criteriaID).Where("market_id", markets).Where("current_price <= ?", maxPrice).Where("current_price >= ? ", minPrice)
//}

func (r AdsRepository) GetAdsForCriteriaPaginated(pagination *Pagination, criteriaID uint, markets []string, minKm *int, maxKm *int, minPrice *int, maxPrice *int) (*[]adsdb.Ad, *Pagination) {
	var ads []adsdb.Ad
	//tx := r.db.Scopes(paginate(ads, pagination, r.db)).Preload("Prices").Preload("Market").Where("criteria_id = ?", criteriaID).Where("market_id", markets).Where("current_price <= ?", maxPrice).Where("current_price >= ? ", minPrice).Find(&ads)
	//tx := r.db.Scopes(withLimitAndOffset(pagination, r.db)).Preload("Prices").Preload("Market").Where("criteria_id = ?", criteriaID).Where("market_id", markets).Where("current_price <= ?", maxPrice).Where("current_price >= ? ", minPrice).Find(&ads)
	tx := r.db.Debug().Scopes(WithFilterConditions(r.db, criteriaID, markets, minKm, maxKm, minPrice, maxPrice)).Find(&ads)
	if tx.Error != nil {
		panic(tx.Error)
	}
	var totalRows int64
	tx = r.db.Debug().Model(&adsdb.Ad{}).Scopes(WithFilterConditions(r.db, criteriaID, markets, minKm, maxKm, minPrice, maxPrice)).Count(&totalRows)
	if tx.Error != nil {
		panic(tx.Error)
	}
	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return &ads, pagination
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

func (r AdsRepository) UpdateCurrentPrice(adID uint) {
	var ad adsdb.Ad
	tx := r.db.Preload("Prices").Find(&ad, adID)
	if tx.Error != nil {
		panic(tx.Error)
	}
	lastPrice := ad.Prices[len(ad.Prices)-1].Price
	ad.CurrentPrice = &lastPrice
	tx = r.db.Save(&ad)
	if tx.Error != nil {
		panic(tx.Error)
	}
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}
func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func AmountGreaterThan1000(db *gorm.DB, maxPrice int) *gorm.DB {
	return db.Where("current_price < ?", maxPrice)
}

func paginate(value interface{}, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Debug().Model(value).Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())
	}
}
