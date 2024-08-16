package properties

type Property struct {
	ID    uint
	Name  string
	Value string
}

type AutoMallProperty struct {
	*Property
	MarketProperties []MarketProperty
}

type MarketProperty struct {
	*Property
	MarketID           uint
	Market             Market
	AutoMallPropertyID uint
	AutoMallProperty   AutoMallProperty
}

type Market struct {
	ID               uint
	Name             string
	MarketProperties []MarketProperty
}
