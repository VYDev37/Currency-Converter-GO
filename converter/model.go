package converter

type Converter struct {
	From        string             `json:"base"`
	Date        string             `json:"date"`
	LastUpdated int64              `json:"time_last_updated`
	Rates       map[string]float64 `json:"rates"`
}
