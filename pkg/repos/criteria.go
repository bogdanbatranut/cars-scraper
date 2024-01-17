package repos

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/valueobjects"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ICriteriaRepository interface {
	GetAll() *[]adsdb.Criteria
	UpdateSelectedCriterias(criterias []valueobjects.Selectable) error
	UpdateSelectedMarkets(markets []valueobjects.Selectable) error
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

func (repo SQLCriteriaRepository) UpdateSelectedCriterias(criterias []valueobjects.Selectable) error {
	transactionErr := repo.db.Transaction(func(tx *gorm.DB) error {
		for _, selectable := range criterias {
			criteria := adsdb.Criteria{}
			repo.db.First(&criteria, selectable.Id)

			criteria.AllowProcess = selectable.Checked

			tx := repo.db.Save(&criteria)
			if tx.Error != nil {
				return tx.Error
			}
		}

		return nil
	})
	if transactionErr != nil {
		return transactionErr
	}
	return nil
}

func (repo SQLCriteriaRepository) UpdateSelectedMarkets(markets []valueobjects.Selectable) error {
	transactionErr := repo.db.Transaction(func(tx *gorm.DB) error {
		for _, selectable := range markets {
			market := adsdb.Market{}
			tx := repo.db.First(&market, selectable.Id)
			if tx.Error != nil {
				return tx.Error
			}

			market.AllowProcess = selectable.Checked

			tx = repo.db.Save(&market)
			if tx.Error != nil {
				return tx.Error
			}
		}

		return nil
	})
	if transactionErr != nil {
		return transactionErr
	}
	return nil
}
