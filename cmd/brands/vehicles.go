package main

import (
	"carscraper/pkg/amconfig"
	"carscraper/pkg/vehicles"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg, err := amconfig.NewViperConfig()
	if err != nil {
		panic(err)
	}

	db, err := inigDB(cfg)
	err = db.AutoMigrate(&vehicles.Vehicle{}, &vehicles.Brand{}, &vehicles.Model{})
	if err != nil {
		panic(err)
	}

	//createSampleData(db)

	//db.Use(dbresolver.Register(dbresolver.Config{
	//	Sources:           nil,
	//	Replicas:          nil,
	//	Policy:            nil,
	//	TraceResolverMode: false,
	//}))
}

func createSampleData(db *gorm.DB) {
	subModel := vehicles.SubModel{
		Name: "GLC-Coupe",
	}

	model := vehicles.Model{
		Name: "GLC",
	}

	market1 := vehicles.Market{
		Name: "autovit",
	}

	brand := vehicles.Brand{
		Name:              "Mercedes-Benz",
		SupportingMarkets: []*vehicles.Market{&market1},
	}

	market2 := vehicles.Market{
		Name: "oferte_bmw",
	}

	vehicle := vehicles.Vehicle{
		Brand:         brand,
		BrandLabel:    brand.Name,
		Model:         model,
		ModelLabel:    model.Name,
		SubModel:      subModel,
		SubModelLabel: subModel.Name,
	}

	//db.Create(&subModel)
	//db.Create(&model)
	db.Create(&market2)
	db.Create(&vehicle)
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

func initAppDB(cfg amconfig.IConfig) (*gorm.DB, error) {
	dsn := createAppDsn(cfg)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	return db, err
}

func createAppDsn(cfg amconfig.IConfig) string {
	databaseName := cfg.GetString(amconfig.AppDBName)
	databaseHost := cfg.GetString(amconfig.AppDBHost)
	dbUser := cfg.GetString(amconfig.AppDBUser)
	dbPass := cfg.GetString(amconfig.AppDBPass)
	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, databaseHost, databaseName)
}
