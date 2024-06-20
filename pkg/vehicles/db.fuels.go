package vehicles

type Fuel struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}
