package repos

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ISellerRepository interface {
	GetAll() (*[]adsdb.Seller, error)
	GetByURLInMarket(string) (*adsdb.Seller, error)
	AddSeller(seller adsdb.Seller) (*adsdb.Seller, error)
}

type SellerRepository struct {
	db *gorm.DB
}

func NewSellerRepository(cfg amconfig.IConfig) *SellerRepository {
	databaseName := cfg.GetString(amconfig.AppDBName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &SellerRepository{
		db: db,
	}
}

func (r SellerRepository) GetAll() (*[]adsdb.Seller, error) {
	var sellers []adsdb.Seller
	tx := r.db.Find(&sellers)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &sellers, nil
}

func (r SellerRepository) GetByURLInMarket(url string) (*adsdb.Seller, error) {
	var seller adsdb.Seller
	tx := r.db.Model(&seller).Where("url_in_market = ?", url).First(&seller)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &seller, nil
}

func (r SellerRepository) AddSeller(seller adsdb.Seller) (*adsdb.Seller, error) {
	tx := r.db.Create(&seller)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &seller, nil
}
