package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

type ModelData struct {
	// Array Input: 0=From, 1=To, 2=Amount
	Inputs     []textinput.Model
	FocusIndex int

	Amount  float64
	Result  float64
	Rate    float64
	Loading bool
	Err     error
}

type Currency struct {
	From        string             `json:"base"`
	Date        string             `json:"date"`
	LastUpdated int64              `json:"time_last_updated`
	Rates       map[string]float64 `json:"rates"`
}

func FetchCurrency(from string) (*Currency, error) {
	res, err := http.Get(fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", from))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data Currency
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	data.LastUpdated = time.Now().Unix()

	cache, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile("cache.json", cache, 0644)

	return &data, nil
}

func LoadCache() (*Currency, error) {
	file, err := os.ReadFile("cache.json")
	if err != nil {
		return nil, err
	}

	var data Currency
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func GetRates(from string) (*Currency, error) {
	data, err := LoadCache()
	if err == nil && data != nil {
		lastUpdated := time.Unix(data.LastUpdated, 0)

		if data.From == from && time.Since(lastUpdated) > time.Minute*30 {
			return data, nil
		}
	}

	return FetchCurrency(from)
}

func GetRate(from, to string) (float64, error) {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	data, err := GetRates(from)
	if err != nil {
		return 0, err
	}

	if data.From == "" {
		return 0, fmt.Errorf("currency with name %s can't be found in database", from)
	}

	rate, ok := data.Rates[to]
	if !ok {
		return 0, fmt.Errorf("currency with name %s can't be found in database", to)
	}

	return rate, nil
}

func Convert(from, to string, amount float64) (*ModelData, error) {
	rate, err := GetRate(from, to)
	if err != nil {
		return nil, err
	}

	return &ModelData{
		Amount: amount,
		Rate:   rate,
		Result: amount * rate,
	}, nil
}
