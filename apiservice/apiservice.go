package apiservice

import (
	"bitcoin/apiserver"
	"bitcoin/bitcoin"
	"bitcoin/conf"
	"context"
	"net/http"

	"github.com/eliona-smart-building-assistant/go-utils/log"
)

// ApiService implements the api service interface in the server pkg
// This service should implement the business logic for every endpoint for the bitcoin api.
type ApiService struct {
}

// NewApiService creates a default api service
func NewApiService() apiserver.ApiService {
	return &ApiService{}
}

// CollectBitcoinRates collects the existing bitcoin rates for the different
// currencies in the defined provider.
func (s *ApiService) GetCurrencyRates(ctx context.Context) (*apiserver.ApiResponse, error) {
	currencyRatesByCode, err := bitcoin.GetCurrencyRates(conf.Endpoint())
	if err != nil {
		log.Error("Bitcoin", "Error during requesting API endpoint: %v", err)
		return &apiserver.ApiResponse{Code: http.StatusInternalServerError}, err
	}

	currencyRates := make([]interface{}, 0)

	for _, cr := range currencyRatesByCode {
		currencyRates = append(currencyRates, cr)
	}

	log.Info("Bitcoin", "got %d currency rates", len(currencyRatesByCode))
	return &apiserver.ApiResponse{Code: http.StatusOK, Body: currencyRates}, nil
}
