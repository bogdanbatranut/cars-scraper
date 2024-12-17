package properties

type AutoMallPropertyKey struct {
	ID               uint   `gorm:"primaryKey"`
	Name             string `gorm:"index"`
	Value            string `gorm:"index"`
	MarketProperties []MarketProperty
}

type MarketProperty struct {
	ID                    uint `gorm:"primaryKey"`
	Value                 string
	MarketID              uint
	AutoMallPropertyKeyID uint
	AutoMallPropertyKey   AutoMallPropertyKey
}

type Market struct {
	ID   uint
	Name string
}
