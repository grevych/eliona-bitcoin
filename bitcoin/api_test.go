package bitcoin

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrencyRates(t *testing.T) {
	updatedTime, _ := time.Parse(coindeskTimeFormat, "2022-09-25T18:45:00+00:00")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"time": {
			"updated": "Sep 25, 2022 18:45:00 UTC",
			"updatedISO": "2022-09-25T18:45:00+00:00",
			"updateduk": "Sep 25, 2022 at 19:45 BST"
			},
			"disclaimer": "This data was produced from the CoinDesk Bitcoin Price Index (USD). Non-USD currency data converted using hourly conversion rate from openexchangerates.org",
			"chartName": "Bitcoin",
			"bpi": {
			"USD": {
			"code": "USD",
			"symbol": "&#36;",
			"rate": "18,879.8173",
			"description": "United States Dollar",
			"rate_float": 18879.8173
			},
			"GBP": {
			"code": "GBP",
			"symbol": "&pound;",
			"rate": "15,775.8243",
			"description": "British Pound Sterling",
			"rate_float": 15775.8243
			},
			"EUR": {
			"code": "EUR",
			"symbol": "&euro;",
			"rate": "18,391.6985",
			"description": "Euro",
			"rate_float": 18391.6985
			}
			}
			}`))
	}))

	currencyRates, err := GetCurrencyRates(server.URL)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(currencyRates))
	assert.Contains(t, currencyRates, &CurrencyRate{
		Code:        "USD",
		Description: "United States Dollar",
		Rate:        18879.8173,
		UpdatedTime: &updatedTime,
	})
	assert.Contains(t, currencyRates, &CurrencyRate{
		Code:        "EUR",
		Description: "Euro",
		Rate:        18391.6985,
		UpdatedTime: &updatedTime,
	})
	assert.Contains(t, currencyRates, &CurrencyRate{
		Code:        "GBP",
		Description: "British Pound Sterling",
		Rate:        15775.8243,
		UpdatedTime: &updatedTime,
	})
}

func TestGetCurrencyRates_WithNoCurrencies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"time": {
			"updated": "Sep 25, 2022 18:45:00 UTC",
			"updatedISO": "2022-09-25T18:45:00+00:00",
			"updateduk": "Sep 25, 2022 at 19:45 BST"
			},
			"disclaimer": "This data was produced from the CoinDesk Bitcoin Price Index (USD). Non-USD currency data converted using hourly conversion rate from openexchangerates.org",
			"chartName": "Bitcoin",
			"bpi": {
			}
			}`))
	}))

	currencyRates, err := GetCurrencyRates(server.URL)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(currencyRates))
}

func TestGetCurrencyRates_WithInvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not what we expected`))
	}))

	currencyRates, err := GetCurrencyRates(server.URL)

	assert.ErrorContains(t, err, "invalid character")
	assert.Equal(t, 0, len(currencyRates))
}

func TestGetCurrencyRates_WithErrorInResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	currencyRates, err := GetCurrencyRates(server.URL)

	assert.ErrorContains(t, err, "error request code")
	assert.Equal(t, 0, len(currencyRates))
}
