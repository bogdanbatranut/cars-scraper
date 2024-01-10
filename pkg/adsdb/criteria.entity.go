package adsdb

import (
	"gorm.io/gorm"
)

type Criteria struct {
	gorm.Model
	Brand        string `json:"brand"`
	CarModel     string `json:"carModel"`
	YearFrom     *int
	YearTo       *int
	Fuel         string
	KmFrom       *int
	KmTo         *int
	AllowProcess bool
	Markets      *[]Market `gorm:"many2many:criteria_markets;"`
}
