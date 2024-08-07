package repos

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/properties"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type PropertiesRepository struct {
	db *gorm.DB
}

func NewPropertiesRepository(cfg amconfig.IConfig) *PropertiesRepository {
	db, err := inigDB(cfg)
	if err != nil {
		panic(err)
	}

	return &PropertiesRepository{db: db}
}

func (r PropertiesRepository) CreateAutoMallProperty(name string, value string) properties.AutoMallProperty {
	p := properties.Property{
		Name:  name,
		Value: value,
	}
	ap := properties.AutoMallProperty{Property: &p}
	r.db.Create(&ap)
	return ap
}

func (r PropertiesRepository) GetMarketProperties(ap properties.AutoMallProperty) []properties.MarketProperty {
	var mp []properties.MarketProperty

	r.db.Preload("Market").Preload("AutoMallProperty").Where(&properties.MarketProperty{
		AutoMallPropertyID: ap.ID,
	}).Find(&mp)
	return mp
}

func (r PropertiesRepository) CreateMarketPropertiesForAutoMallPropertiesInMarket(string, name string, value string, ap properties.AutoMallProperty, marketID uint) {
	var market properties.Market
	r.db.First(&market, marketID)
	mp := properties.CreateMarketProperty(name, value)
	mp.Market = market
	mp.AutoMallProperty = ap
	r.db.Create(&mp)
}

func (r PropertiesRepository) GetAutoMallPropertyByName(name string) *[]properties.AutoMallProperty {
	var ap []properties.AutoMallProperty
	r.db.Preload("MarketProperties").Where(&properties.AutoMallProperty{
		Property: &properties.Property{
			Name: name,
		},
	}).Find(&ap)
	return &ap
}

func (r PropertiesRepository) GetAutoMallPropertyByID(id uint) *properties.AutoMallProperty {
	var ap properties.AutoMallProperty
	r.db.Preload("MarketProperties").Where(&properties.AutoMallProperty{
		Property: &properties.Property{
			ID: id,
		},
	}).First(&ap)
	return &ap
}

func inigDB(cfg amconfig.IConfig) (*gorm.DB, error) {
	dsn := createMapperDsn(cfg)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	return db, err
}

func createMapperDsn(cfg amconfig.IConfig) string {
	databaseName := cfg.GetString(amconfig.AppDBMapperName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
}
