package repos

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IMarketsRepository interface {
	GetAll() *[]adsdb.Market
}

type SQLMarketsRepository struct {
	db *gorm.DB
}

func NewSQLMarketsRepository(cfg amconfig.IConfig) *SQLMarketsRepository {
	databaseName := cfg.GetString(amconfig.AppDBName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &SQLMarketsRepository{
		db: db,
	}
}

func (repo SQLMarketsRepository) GetAll() *[]adsdb.Market {
	var markets []adsdb.Market
	repo.db.Preload("Criterias").Find(&markets)
	return &markets
}
