package main

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"carscraper/pkg/errorshandler"
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := amconfig.NewViperConfig()
	errorshandler.HandleErr(err)

	databaseName := cfg.GetString(amconfig.AppDBName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	var ads []adsdb.Ad
	tx := db.Where("ad_url NOT LIKE ?", "https%").Find(&ads)
	if tx.Error != nil {
		panic(tx.Error)
	}

	log.Printf("Found %d ads", len(ads))

	for index, foundAd := range ads {
		log.Printf("Changing ad number %d of %d", index, len(ads))
		if foundAd.MarketID == 11 {
			if !strings.Contains(foundAd.Ad_url, "www.mobile.de") {
				tx := db.Model(&foundAd).Update("ad_url", fmt.Sprintf("https://www.mobile.de%s", foundAd.Ad_url))
				if tx.Error != nil {
					err = tx.Error
				}
			}

		}

		if foundAd.MarketID == 12 {
			if !strings.Contains(foundAd.Ad_url, "www.autoscout24.ro") {
				tx := db.Model(&foundAd).Update("ad_url", fmt.Sprintf("https://www.autoscout24.ro%s", foundAd.Ad_url))
				if tx.Error != nil {
					err = tx.Error
				}
			}
		}
	}

}
