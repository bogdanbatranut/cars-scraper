package main

import (
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := amconfig.NewViperConfig()
	if err != nil {
		panic(err)
	}

	db, err := initAppDB(cfg)
	migrate(db)
	if err != nil {
		panic(err)
	}
}

func migrate(db *gorm.DB) {
	err := db.Debug().AutoMigrate(&adsdb.SessionLog{}, &adsdb.CriteriaLog{}, &adsdb.PageLog{})
	if err != nil {
		panic(err)
	}
}

func initAppDB(cfg amconfig.IConfig) (*gorm.DB, error) {
	dsn := createAppDsn(cfg)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	return db, err
}

func createAppDsn(cfg amconfig.IConfig) string {
	databaseName := cfg.GetString(amconfig.AppDBLogsName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
}
