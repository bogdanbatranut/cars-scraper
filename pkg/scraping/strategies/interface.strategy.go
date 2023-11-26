package strategies

type Ad struct {
	Brand      string
	Model      string
	Year       int
	Km         int
	Fuel       string
	Price      int
	AdID       *int
	Ad_url     string
	SellerType string
	SellerName *string
	SellerURL  *string
}

type IScrapingStrategy interface {
	Execute(url string) ([]Ad, error)
}
