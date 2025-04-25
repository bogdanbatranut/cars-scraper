package repos

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/properties"
	"errors"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	Brand                = "brand"
	Model                = "model"
	SubModel             = "subModel"
	PathBrand            = "pathBrand"
	PathModel            = "pathModel"
	QueryParamBrandParam = "queryParamBrand"
	QueryParamModelParam = "queryParamModel"
	QueryParamFuelParam  = "queryParamFuel"
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

func (r PropertiesRepository) Migrate() {
	err := r.db.Debug().AutoMigrate(&properties.AutoMallPropertyKey{}, &properties.Market{}, &properties.MarketProperty{})
	if err != nil {
		panic(err)
	}
}

func (r PropertiesRepository) CreateAutoMallPropertyKey(name string, value string) properties.AutoMallPropertyKey {
	p := properties.AutoMallPropertyKey{
		Name:  name,
		Value: value,
	}
	r.db.Create(&p)
	return p
}

func (r PropertiesRepository) GetAllMarketProperties(ap properties.AutoMallPropertyKey) []properties.MarketProperty {
	var mp []properties.MarketProperty

	r.db.Preload("Market").Preload("AutoMallProperty").Where(&properties.MarketProperty{
		AutoMallPropertyKeyID: ap.ID,
	}).Find(&mp)
	return mp
}

//func (r PropertiesRepository) GetAutoMallPropertyByName(name string) *[]properties.AutoMallProperty {
//	var ap []properties.AutoMallProperty
//	r.db.Preload("MarketProperties").Where(&properties.AutoMallProperty{
//		Property: &properties.Property{
//			Name: name,
//		},
//	}).Find(&ap)
//	return &ap
//}

//func (r PropertiesRepository) GetAutoMallPropertyByID(id uint) *properties.AutoMallProperty {
//	var ap properties.AutoMallProperty
//	r.db.Preload("MarketProperties.Market").Where(&properties.AutoMallProperty{
//		Property: &properties.Property{
//			ID: id,
//		},
//	}).First(&ap)
//	return &ap
//}

func (r PropertiesRepository) GetPropertyMarketValues(id uint, marketID uint) *properties.MarketProperty {
	var ap properties.MarketProperty
	r.db.Debug().Preload("AutoMallPropertyKey").Where(&properties.MarketProperty{
		MarketID:              marketID,
		AutoMallPropertyKeyID: id,
	}).First(&ap)
	return &ap
}

func (r PropertiesRepository) GetPropertyMarketValuesForTypeAndValue(propName string, value string, marketID uint) *properties.MarketProperty {
	var ap properties.MarketProperty

	r.db.Debug().Preload("AutoMallPropertyKey").Where("market_id", marketID).
		Where("auto_mall_property_key_id = (?)", r.db.Table("auto_mall_property_keys").Select("id").Where(properties.AutoMallPropertyKey{
			Name:  propName,
			Value: value,
		}).Limit(1)).First(&ap)
	return &ap
}

func (r PropertiesRepository) addMarketProperty(property properties.MarketProperty) {
	var existingProperty properties.MarketProperty
	r.db.Debug().Where(property).First(&existingProperty)

	if existingProperty.ID == 0 {
		r.db.Debug().Omit(clause.Associations).Create(&property)
	}
}

func (r PropertiesRepository) Test(automallPropertyKey properties.AutoMallPropertyKey) {
	err := r.db.Debug().Where(automallPropertyKey).Limit(1).First(&automallPropertyKey).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("not found....")
	}
}

func (r PropertiesRepository) AddMarketProperties(automallPropertyKey properties.AutoMallPropertyKey) {
	var existingAutoMallPropertyKey properties.AutoMallPropertyKey

	r.db.Debug().Where(automallPropertyKey).Limit(1).First(&existingAutoMallPropertyKey)
	if existingAutoMallPropertyKey.ID == 0 {
		r.db.Omit(clause.Associations).Create(&automallPropertyKey)

		for _, property := range automallPropertyKey.MarketProperties {
			property.AutoMallPropertyKeyID = automallPropertyKey.ID
			r.addMarketProperty(property)
		}
	} else {
		for _, property := range automallPropertyKey.MarketProperties {
			property.AutoMallPropertyKeyID = existingAutoMallPropertyKey.ID
			r.addMarketProperty(property)
		}
	}

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
