package repos

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IAdsRepository interface {
	GetAll() (*[]adsdb.Ad, error)
	Upsert(ad adsdb.Ad) error
	//UpsertAll(ads []adsdb.Ad) error
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

func (r AdsRepository) Upsert(ad adsdb.Ad) error {
	tx := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&ad)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

//func (r AdsRepository) UpserAll(ads *adsdb.Ad) error {
//	transactionErr := r.db.Transaction(func(tx *gorm.DB) error {
//		var err error
//		for _, ad := range ads {
//			tx := r.db.Create(&ad)
//			if tx.Error != nil {
//				err = tx.Error
//			}
//		}
//		if err != nil {
//			return err
//		}
//		return nil
//	})
//	return transactionErr
//}
