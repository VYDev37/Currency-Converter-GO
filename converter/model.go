package converter

type Converter struct {
	From        string             `json:"base"`
	Date        string             `json:"date"`
	LastUpdated int64              `json:"time_last_updated`
	Rates       map[string]float64 `json:"rates"`
}

type ModelData struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Rate   float64 `json:"rate"`
	Amount float64 `json:"amount"`
	Result float64 `json:"result"`
}
