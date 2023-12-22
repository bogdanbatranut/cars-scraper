package results

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ResultsRepository struct {
	db *gorm.DB
}

func NewResultsRepository(cfg amconfig.IConfig) *ResultsRepository {
	databaseName := cfg.GetString(amconfig.AppDBName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &ResultsRepository{
		db: db,
	}
}

// WriteResults writes results in a transaction
// if it fails returns an error
func (r ResultsRepository) WriteResults(ads []adsdb.Ad) error {
	r.db.Session(&gorm.Session{FullSaveAssociations: true})
	transactionErr := r.db.Transaction(func(tx *gorm.DB) error {
		var err error
		for _, ad := range ads {
			price := ad.Prices[0]
			muuid := ad.MarketUUID
			tx = r.db.FirstOrCreate(&ad, adsdb.Ad{MarketUUID: muuid}, adsdb.Ad{Prices: ad.Prices}, adsdb.Ad{SellerID: ad.SellerID})
			if tx.Error != nil {
				err = tx.Error
			}

			price.AdID = ad.ID
			tx = r.db.FirstOrCreate(&price, adsdb.Price{Price: price.Price})
			if tx.Error != nil {
				err = tx.Error
			}
		}
		if err != nil {
			return err
		}
		return nil
	})
	return transactionErr
}

// WriteResults writes results in a transaction
// if it fails returns an error
func (r ResultsRepository) Write(ads []adsdb.Ad) (*[]uint, error) {
	adsIds := []uint{}
	r.db.Session(&gorm.Session{FullSaveAssociations: true})
	transactionErr := r.db.Transaction(func(tx *gorm.DB) error {
		var err error
		for _, ad := range ads {
			price := ad.Prices[0]
			muuid := ad.MarketUUID
			tx = r.db.FirstOrCreate(&ad, adsdb.Ad{MarketUUID: muuid}, adsdb.Ad{Prices: ad.Prices}, adsdb.Ad{SellerID: ad.SellerID})
			if tx.Error != nil {
				err = tx.Error
			}

			price.AdID = ad.ID
			tx = r.db.FirstOrCreate(&price, adsdb.Price{Price: price.Price})
			if tx.Error != nil {
				err = tx.Error
			}
			adsIds = append(adsIds, ad.ID)
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

func (r ResultsRepository) GetAllAdsIDs(marketID uint, criteriaID uint) *[]uint {
	var ads []adsdb.Ad
	var adsIDs []uint
	r.db.Model(adsdb.Ad{}).Select("id").Where("market_id = ? AND criteria_id = ?", marketID, criteriaID).Find(&ads)
	for _, ad := range ads {
		adsIDs = append(adsIDs, ad.ID)
	}
	return &adsIDs
}

func (r ResultsRepository) DeleteAd(adID uint) {
	var ad adsdb.Ad
	r.db.Model(adsdb.Ad{}).Delete(&ad, adID)
}
