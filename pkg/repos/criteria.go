package repos

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CriteriaRepository interface {
	GetAll() *[]adsdb.Criteria
}

type SQLCriteriaRepository struct {
	db *gorm.DB
}

func NewSQLCriteriaRepository(cfg config.IConfig) *SQLCriteriaRepository {
	databaseName := cfg.GetString(config.AppDBName)
	databaseHost := cfg.GetString(config.AppDBHost)
	dbUser := cfg.GetString(config.AppDBUser)
	dbPass := cfg.GetString(config.AppDBPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &SQLCriteriaRepository{
		db: db,
	}

	//databaseName := "carsfinder"
	//dsn := fmt.Sprintf("root:rootpass@tcp(host.docker.internal:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", databaseName)
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//if err != nil {
	//	panic(err)
	//}
	//return &SQLCriteriaRepository{
	//	db: db,
	//}
}

func (repo SQLCriteriaRepository) GetAll() *[]adsdb.Criteria {
	var criterias []adsdb.Criteria
	repo.db.Preload("Markets").Find(&criterias)
	return &criterias
}
