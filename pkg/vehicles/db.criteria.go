package vehicles

type Criteria struct {
	ID            uint `gorm:"primaryKey"`
	BrandID       uint `gorm:"foreignKey:BrandID, references:ID"`
	BrandLabel    string
	CarModelID    uint `gorm:"foreignKey:CarModelID, references:ID"`
	CarModelLabel string
	YearFrom      *int
	YearTo        *int
	Fuel          Fuel
	FuelLabel     string
	KmFrom        *int
	KmTo          *int
}
