package initialization

import (
	"carscraper/pkg/amconfig"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitVehiclesDB() (*gorm.DB, error) {
	cfg, err := amconfig.NewViperConfig()
	if err != nil {
		return nil, err
	}
	return inigDB(cfg)
}

func InitAutoMallDB() (*gorm.DB, error) {
	cfg, err := amconfig.NewViperConfig()
	if err != nil {
		return nil, err
	}
	dsn := createAppDsn(cfg)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	return db, err
}

func inigDB(cfg amconfig.IConfig) (*gorm.DB, error) {
	dsn := createVehiclesDsn(cfg)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	return db, err
}

func createVehiclesDsn(cfg amconfig.IConfig) string {
	databaseName := cfg.GetString(amconfig.AppDBVehiclesName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
}

func createAppDsn(cfg amconfig.IConfig) string {
	databaseName := cfg.GetString(amconfig.AppDBName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
}
