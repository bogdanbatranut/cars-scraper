package repos

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
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

func NewTestMigrationRepository(cfg amconfig.IConfig) *MigrationRepository {
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	dbHost := cfg.GetString(amconfig.AppBaseURL)
	dbName := cfg.GetString(amconfig.AppTestDBName)

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

func NewMigrationRepository(cfg amconfig.IConfig) *MigrationRepository {
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	dbHost := cfg.GetString(amconfig.AppBaseURL)
	dbName := cfg.GetString(amconfig.AppDBName)

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
