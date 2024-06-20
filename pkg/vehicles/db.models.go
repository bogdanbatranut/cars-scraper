package vehicles

type Vehicle struct {
	ID            uint `gorm:"primaryKey"`
	BrandID       uint `gorm:"foreignKey:BrandID, references:ID"`
	BrandLabel    string
	Brand         Brand
	ModelID       uint `gorm:"foreignKey:ModelID, references:ID"`
	Model         Model
	ModelLabel    string
	SubModelID    uint `gorm:"foreignKey:SubModelID, references:ID"`
	SubModel      SubModel
	SubModelLabel string

	//BodyTypeID    uint `gorm:"foreignKey:BodyTypeID, references:ID"`
	//BodyType      BodyType
	//BodyTypeLabel string
}

type Market struct {
	ID              uint `gorm:"primaryKey"`
	Name            string
	SupportedBrands []*Brand `gorm:"many2many:market_brands"`
}

type Brand struct {
	ID                uint `gorm:"primaryKey"`
	Name              string
	Models            []Model
	SupportingMarkets []*Market `gorm:"many2many:market_brands"`
}

type Model struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	BrandID   uint `gorm:"foreignKey:BrandID, references:ID"`
	Brand     Brand
	SubModels []SubModel
}

type SubModel struct {
	ID      uint `gorm:"primaryKey"`
	Name    string
	ModelID uint `gorm:"foreignKey:ModelID, references:ID"`
	Model   Model
}

// TODO leave this for now
//type BodyType struct {
//	ID   uint `gorm:"primaryKey"`
//	Name string
//}
