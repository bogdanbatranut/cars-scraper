package repos

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/config"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IMigrationRepository interface {
	Migrate(...interface{}) error
}

type MigrationRepository struct {
	db *gorm.DB
}

func NewMigrationRepository(cfg config.IConfig) *MigrationRepository {
	dbUser := cfg.GetString(config.AppDBUser)
	dbPass := cfg.GetString(config.AppDBPass)
	dbHost := cfg.GetString(config.AppBaseURL)
	dbName := cfg.GetString(config.AppDBName)

	//dsn := fmt.Sprintf("root:rootpass@tcp(host.docker.internal:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", databaseName)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &MigrationRepository{
		db: db,
	}
}

func (repo MigrationRepository) Migrate(tables ...interface{}) error {
	err := repo.db.AutoMigrate(tables)
	return err
}

func (repo MigrationRepository) GetDB() *gorm.DB {
	return repo.db
}

func (repo MigrationRepository) WriteMarkets(m []adsdb.Market) *[]adsdb.Market {
	repo.db.Create(&m)
	return &m
}

func (repo MigrationRepository) WriteCriterias(markets *[]adsdb.Market, criterias *[]adsdb.Criteria) {
	for _, c := range *criterias {
		c.Markets = markets
		repo.db.Save(&c)
	}
}
