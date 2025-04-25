package autovit

type AutovitPOSTPayload struct {
	Query     string `json:"query"`
	Variables struct {
		Filters []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"filters"`
		ItemsPerPage               int      `json:"itemsPerPage"`
		Page                       int      `json:"page"`
		Parameters                 []string `json:"parameters"`
		ShouldEncryptSensitiveData bool     `json:"shouldEncryptSensitiveData"`
		SortBy                     string   `json:"sortBy"`
	} `json:"variables"`
}
