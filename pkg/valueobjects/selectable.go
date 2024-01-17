package valueobjects

type Selectable struct {
	Id      int    `json:"id"`
	Checked bool   `json:"checked"`
	Label   string `json:"label"`
}
