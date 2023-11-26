package strategies

type ImplementationStrategies struct {
	strategies map[string]IScrapingStrategy
}

func NewImplemetationStrategies() ImplementationStrategies {
	s := make(map[string]IScrapingStrategy)
	// here we add to the map the implementations ...
	s["autovit"] = NewAutovitStrategy()
	s["mobile.de"] = NewMobileDeStrategy()

	is := ImplementationStrategies{
		strategies: s,
	}
	return is
}

func (is ImplementationStrategies) GetImplementation(marketName string) IScrapingStrategy {
	return is.strategies[marketName]
}
