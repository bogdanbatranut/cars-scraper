package properties

func CreateMarketProperty(name string, value string) MarketProperty {
	p := Property{
		Name:  name,
		Value: value,
	}

	mp := MarketProperty{
		Property:           &p,
		MarketID:           0,
		Market:             Market{},
		AutoMallPropertyID: 0,
		AutoMallProperty:   AutoMallProperty{},
	}
	return mp
}
