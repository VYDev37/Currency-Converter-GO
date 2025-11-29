package converter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"
	"sort"
	"strings"
	"time"
)

type ConvFunc interface {
	Init() error
	Convert(string, string, float64) (*Converter, error)
	GetCurrencies() []string
}

type ConvManager struct {
	currencies []string
	rates      map[string]float64
}

func FetchCurrency(from string) (*Converter, error) {
	res, err := http.Get(fmt.Sprintf("https://api.exchangerate-api.com/v4/latest/%s", from))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data Converter
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	data.LastUpdated = time.Now().Unix()

	cache, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile("cache.json", cache, 0644)

	return &data, nil
}

func LoadCache() (*Converter, error) {
	file, err := os.ReadFile("cache.json")
	if err != nil {
		return nil, err
	}

	var data Converter
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

func LoadRates(from string) (*Converter, error) {
	data, err := LoadCache()
	if err == nil && data != nil {
		lastUpdated := time.Unix(data.LastUpdated, 0)

		if time.Since(lastUpdated) < time.Minute*30 {
			return data, nil
		}
	}

	return FetchCurrency(from)
}

// manager functions
func (manager *ConvManager) Init() error {
	data, err := LoadRates("USD") // to get all existing currency ids

	if err != nil {
		return err
	}

	manager.rates = data.Rates
	manager.currencies = make([]string, 0, len(manager.rates))

	for key := range manager.rates {
		manager.currencies = append(manager.currencies, key)
	}

	sort.Strings(manager.currencies)
	return nil
}

func (manager *ConvManager) Convert(from, to string, amount float64) (float64, error) {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	if !slices.Contains(manager.currencies, from) || !slices.Contains(manager.currencies, to) {
		if err := manager.Init(); err != nil { // reload
			return 0, err
		}
	}

	rateFrom := manager.rates[from]
	rateTo := manager.rates[to]

	// set USD as base
	if from == "USD" {
		rateFrom = 1
	}
	if to == "USD" {
		rateTo = 1
	}
	if rateFrom == 0 {
		return 0, fmt.Errorf("currency %s can't be found in database", from)
	}
	if rateTo == 0 {
		return 0, fmt.Errorf("currency %s can't be found in database", to)
	}

	convertedRate := rateTo / rateFrom
	result := convertedRate * amount

	return result, nil
}

func (manager *ConvManager) GetCurrencies() []string {
	return manager.currencies
}

func (manager *ConvManager) GetRates() map[string]float64 {
	return manager.rates
}
