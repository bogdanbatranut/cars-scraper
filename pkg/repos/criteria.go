package repos

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ICriteriaRepository interface {
	GetAll() *[]adsdb.Criteria
}

type SQLCriteriaRepository struct {
	db *gorm.DB
}

func NewSQLCriteriaRepository(cfg amconfig.IConfig) *SQLCriteriaRepository {
	databaseName := cfg.GetString(amconfig.AppDBName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &SQLCriteriaRepository{
		db: db,
	}
}

func (repo SQLCriteriaRepository) GetAll() *[]adsdb.Criteria {
	var criterias []adsdb.Criteria
	repo.db.Model(&adsdb.Criteria{}).Preload("Markets").Order("brand").Order("car_model").Find(&criterias)
	return &criterias
}
