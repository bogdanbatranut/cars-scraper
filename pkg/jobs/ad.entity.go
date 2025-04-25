package jobs

type Ad struct {
	Title              *string
	Brand              string
	Model              string
	Year               int
	Km                 int
	Fuel               string
	Price              int
	AdID               string
	Ad_url             string
	SellerType         string
	SellerName         *string
	SellerNameInMarket *string
	SellerOwnURL       *string
	SellerMarketURL    *string
	Thumbnail          *string
}

func (ad Ad) SetBrand(brand string) {
	ad.Brand = brand
}

func (ad Ad) SetModel(model string) {
	ad.Model = model
}

func (ad Ad) SetFuel(fuel string) {
	ad.Brand = fuel
}
