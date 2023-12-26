package cardb

import (
	"carscraper/pkg/jobs"

	"gorm.io/gorm"
)

type Manufacturer struct {
	*gorm.Model
	Name string
}

type Model struct {
	*gorm.Model
	Name string
}

type SearchTerms struct {
	*gorm.Model
}

type MarketManufacturer struct {
	*gorm.Model
	MarketID       uint
	ManufacturerID uint
	Name           string
	Value          string
	Market         jobs.Market
	MarketModels   []MarketModel
}

type MarketModel struct {
	*gorm.Model
	Name                 string
	Value                string
	MarketManufacturerID uint
	MarketManufacturer   MarketManufacturer
}

type MarketFuel struct {
	*gorm.Model
	Name   string
	Values []MarketFuelValues
}

type MarketFuelValues struct {
	*gorm.Model
	Value  string
	FuelID uint
}
