package repos

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ChartsRepository struct {
	db *gorm.DB
}

func NewChartsRepository(cfg amconfig.IConfig) *ChartsRepository {
	databaseName := cfg.GetString(amconfig.AppDBName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &ChartsRepository{
		db: db,
	}
}

func (r ChartsRepository) GetAdsPricesByStep(step int) {
	// get max price for ad
	//var price adsdb.Price

	var allPrices []adsdb.Price
	//tx := r.db.Debug().Find(&allPrices).Group("ad_id")
	//if tx.Error != nil {
	//	panic(tx.Error)
	//}
	//
	//for i, price := range allPrices {
	//	log.Printf("%d - %d - %d", price.ID, price.AdID, price.Price)
	//
	//	if i == 100 {
	//		return
	//	}
	//}
	type Res struct {
		Max int
	}
	var res Res
	tx := r.db.Debug().Raw("select max(price) as max from prices where id in (select max(id) from prices group by ad_id)").Scan(&res)
	if tx.Error != nil {
		panic(tx.Error)
	}
	log.Printf("%+v", res)
	//maxPrice := res.Max

	tx = r.db.Debug().Raw("select min(price) as max from prices where id in (select max(id) from prices group by ad_id)").Scan(&res)
	if tx.Error != nil {
		panic(tx.Error)
	}
	//minPrice := res.Max

	tx = r.db.Raw("select * from prices where id in (select max(id) from prices group by ad_id)").Scan(&allPrices)
	if tx.Error != nil {
		panic(tx.Error)
	}
}
