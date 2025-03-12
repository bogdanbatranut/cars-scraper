package autovit

type VariablesRequestParamValue struct {
	After                   interface{}  `json:"after"`
	Experiments             []Experiment `json:"experiments"`
	Filters                 []Filters    `json:"filters"`
	IncludeCepik            bool         `json:"includeCepik"`
	IncludeFiltersCounters  bool         `json:"includeFiltersCounters"`
	IncludeNewPromotedAds   bool         `json:"includeNewPromotedAds"`
	IncludePriceEvaluation  bool         `json:"includePriceEvaluation"`
	IncludePromotedAds      bool         `json:"includePromotedAds"`
	IncludeRatings          bool         `json:"includeRatings"`
	IncludeSortOptions      bool         `json:"includeSortOptions"`
	IncludeSuggestedFilters bool         `json:"includeSuggestedFilters"`
	MaxAge                  int          `json:"maxAge"`
	Page                    int          `json:"page"`
	Parameters              []string     `json:"parameters"`
	PromotedInput           struct {
	} `json:"promotedInput"`
	SearchTerms interface{} `json:"searchTerms"`
	SortBy      string      `json:"sortBy"`
}

type ExtensionsRequestParamValue struct {
	PersistedQuery struct {
		Sha256Hash string `json:"sha256Hash"`
		Version    int    `json:"version"`
	} `json:"persistedQuery"`
}
