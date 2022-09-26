package apiserver

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ApiController binds http requests to an api service and writes the service
// results to the http response
type ApiController struct {
	service ApiService
}

// ApiResponse represents a very basic http response
type ApiResponse struct {
	Code int
	Body interface{}
}

// ApiService defines the available methods in the service
type ApiService interface {
	GetCurrencyRates(context.Context) (*ApiResponse, error)
}

// A Route defines the parameters for an api endpoint
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Router defines the required methods for retrieving api routes
type Router interface {
	Routes() []Route
}

// NewApiController creates a default api controller
func NewApiController(s ApiService) Router {
	controller := &ApiController{
		service: s,
	}

	return controller
}

// EncodeJSONResponse uses the json encoder to write an interface to the http
// response with an optional status code
func EncodeJSONResponse(i interface{}, status *int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if status != nil {
		w.WriteHeader(*status)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if err := json.NewEncoder(w).Encode(i); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

// NewRouter creates a new router for any number of api routers
func NewRouter(routers ...Router) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, api := range routers {
		for _, route := range api.Routes() {
			var handler http.Handler
			handler = route.HandlerFunc
			// handler = Logger(handler, route.Name)

			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		}
	}

	return router
}

// Routes returns all the api routes for the ApiController
func (c *ApiController) Routes() []Route {
	return []Route{
		{
			"GetCurrencyRates",
			strings.ToUpper("Get"),
			"/v2/bitcoin/rates",
			c.GetCurrencyRates,
		},
	}
}

// GetCurrencyRates attempts to get currency rates from bitcoin api
func (c *ApiController) GetCurrencyRates(w http.ResponseWriter, r *http.Request) {
	result, err := c.service.GetCurrencyRates(r.Context())
	// If an error occurred, encode the error with the status code
	if err != nil {
		status := http.StatusInternalServerError
		EncodeJSONResponse(err.Error(), &status, w)
		return
	}

	EncodeJSONResponse(result.Body, &result.Code, w)
}
