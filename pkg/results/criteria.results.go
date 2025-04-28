package results

type ProcessedCriteria struct {
	CriteriaID uint
	Markets    *[]ProcessedMarket
}

func NewProcessedCriteria(criteriaID uint) *ProcessedCriteria {
	pm := make([]ProcessedMarket, 0)
	return &ProcessedCriteria{
		CriteriaID: criteriaID,
		Markets:    &pm,
	}
}

func (pc *ProcessedCriteria) getMarket(marketID uint) *ProcessedMarket {
	for _, market := range *pc.Markets {
		if market.MarketID == marketID {
			return &market
		}
	}
	market := NewProcessedMarket(marketID)
	*pc.Markets = append(*pc.Markets, *market)
	return market
}

func (pc ProcessedCriteria) isComplete() bool {
	if pc.Markets == nil {
		return false
	}
	for _, market := range *pc.Markets {
		if !market.isComplete() {
			return false
		}
	}
	return true
}
