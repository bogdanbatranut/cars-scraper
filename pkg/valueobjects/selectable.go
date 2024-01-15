package valueobjects

type Selectable struct {
	Id       int    `json:"id"`
	Selected bool   `json:"selected"`
	Label    string `json:"label"`
}
