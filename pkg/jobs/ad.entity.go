package jobs

type Ad struct {
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
