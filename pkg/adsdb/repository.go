package adsdb

import (
	"carscraper/pkg/errorshandler"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBRepo struct {
	db *gorm.DB
}

func NewRepository(databaseName string) DBRepo {
	dsn := fmt.Sprintf("root:rootpass@tcp(host.docker.internal:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	errorshandler.HandleErr(err)

	return DBRepo{db: db}
}

func (repo DBRepo) Migrate() error {
	return repo.db.Migrator().AutoMigrate(
		&Seller{}, &Criteria{},
		&Market{}, &Ad{}, &Price{},
		&SellerMarkets{}, &CriteriaMarkets{})
}

func (repo DBRepo) AddCriteria(c Criteria) {
	repo.db.Save(&c)
}
