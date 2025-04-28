package results

type MarketResults struct {
	MarketID     uint
	ResultsPages MarketScrapingResults
}
type ProcessedMarket struct {
	MarketID uint
	Results  *[]PageResult
}

func NewProcessedMarket(marketID uint) *ProcessedMarket {
	pr := make([]PageResult, 0)
	return &ProcessedMarket{
		MarketID: marketID,
		Results:  &pr,
	}
}

func (pm *ProcessedMarket) AddPageResult(pageResults PageResult) {
	*pm.Results = append(*pm.Results, pageResults)
}

func (pm ProcessedMarket) isComplete() bool {
	hasLastPage := false
	lastPageNumber := 0
	existingPageNumbers := make(map[int]bool)

	for _, page := range *pm.Results {
		existingPageNumbers[page.pageNumber] = true
		if page.isLastPage {
			hasLastPage = true
			lastPageNumber = page.pageNumber
		}
	}
	if !hasLastPage {
		return false
	}

	// has all pages
	for i := 1; i <= lastPageNumber; i++ {
		if _, exists := existingPageNumbers[i]; !exists {
			return false
		}
	}

	return true
}
