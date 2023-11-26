package urlbuilder

type PageURL struct {
	MarketName string
	MarketURL  string
	CarBrand   string
	CarModel   string
	YearFrom   *int
	YearTo     *int
	Fuel       *string
	KmFrom     *int
	KmTo       *int
	PageNumber int
}
