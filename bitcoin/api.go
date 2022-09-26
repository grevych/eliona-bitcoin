package bitcoin

import (
	"encoding/json"
	"time"

	"github.com/eliona-smart-building-assistant/go-utils/http"
)

// CoindeskResponse represents the response coming from coindesk endpoint
// when quering the available currency rates for bitcoin
type CoindeskResponse struct {
	Time *struct {
		Updated string `json:"updatedISO"`
	} `json:"time"`
	BPI map[string]*struct {
		Code        string  `json:"code"`
		Description string  `json:"description"`
		Rate        float64 `json:"rate_float"`
	} `json:"bpi"`
}

// CurrencyRate represents the business logic model for each bitcoin currency rate
type CurrencyRate struct {
	Code        string     `json:"code"`
	Description string     `json:"description"`
	Rate        float64    `json:"rate"`
	UpdatedTime *time.Time `json:"updated_time"`
}

const coindeskTimeFormat = "2006-01-02T15:04:05-07:00"

// GetCurrencyRates returns the current available currency rates for bitcoin
// using Coindesk as the provider
func GetCurrencyRates(ep string) ([]*CurrencyRate, error) {
	var (
		currencyRates []*CurrencyRate
		updatedTime   time.Time
	)

	payload, err := request(ep)
	if err != nil {
		return nil, err
	}

	var result CoindeskResponse
	err = json.Unmarshal(payload, &result)
	if err != nil {
		return nil, err
	}

	updatedTime, err = time.Parse(coindeskTimeFormat, result.Time.Updated)
	if err != nil {
		return nil, err
	}

	for _, v := range result.BPI {
		currencyRates = append(currencyRates, &CurrencyRate{
			v.Code, v.Description, v.Rate, &updatedTime,
		})
	}

	return currencyRates, nil
}

// request calls the provider's api to get the bitcoin rates for the available currencies
func request(ep string) ([]byte, error) {
	request, err := http.NewRequest(ep)
	if err != nil {
		return nil, err
	}
	payload, err := http.Do(request, time.Second*10, true)
	return payload, err
}
