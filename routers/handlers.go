package routers

import (
	"corona/helpers"
	"corona/middleware"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type (
	HandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, *helpers.Error)
)

func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var errs []string
	r.ParseForm()
	data, err := fn(w, r)
	if err != nil {
		errs = append(errs, err.Error())
		w.WriteHeader(err.StatusCode)
	}
	resp := helpers.Response{
		Data: data,
		BaseResponse: helpers.BaseResponse{
			Errors: errs,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		return
	}
}

func InitHandlers() *mux.Router {
	r := mux.NewRouter()

	http.Handle("/", r)

	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	apiV1.Handle("/coronavirus", middleware.TokenMiddleware(middleware.RateLimitMiddleware(
		HandlerFunc(HandlerCoronaList)))).Methods(http.MethodGet)
	apiV1.Handle("/coronavirus/continents/{continent}", middleware.TokenMiddleware(middleware.RateLimitMiddleware(
		HandlerFunc(HandlerCoronaByContinent)))).Methods(http.MethodGet)
	apiV1.Handle("/coronavirus/countries/{country}", middleware.TokenMiddleware(middleware.RateLimitMiddleware(
		HandlerFunc(HandlerCoronaByCountry)))).Methods(http.MethodGet)

	apiV1.Handle("/coronavirus",
		HandlerFunc(HandlerCoronaAdd)).Methods(http.MethodPost)

	apiV1.Handle("/countries",
		HandlerFunc(HandlerCountryList)).Methods(http.MethodGet)
	apiV1.Handle("/countries",
		HandlerFunc(HandlerCountryAdd)).Methods(http.MethodPost)
	apiV1.Handle("/countries/{id}",
		HandlerFunc(HandlerCountryDetail)).Methods(http.MethodGet)

	apiV1.Handle("/continents",
		HandlerFunc(HandlerContinentList)).Methods(http.MethodGet)
	apiV1.Handle("/continents/{id}",
		HandlerFunc(HandlerContinentDetail)).Methods(http.MethodGet)

	apiV1.Handle("/subscriptions",
		HandlerFunc(HandlerSubscriptionList)).Methods(http.MethodGet)
	apiV1.Handle("/subscriptions/{id}",
		HandlerFunc(HandlerSubscriptionDetail)).Methods(http.MethodGet)

	apiV1.Handle("/rate_limits",
		HandlerFunc(HandlerRateLimitList)).Methods(http.MethodGet)
	apiV1.Handle("/rate_limits/{id}",
		HandlerFunc(HandlerRateLimitDetail)).Methods(http.MethodGet)

	apiV1.Handle("/tokens",
		HandlerFunc(HandlerTokenList)).Methods(http.MethodGet)
	apiV1.Handle("/tokens/{id}",
		HandlerFunc(HandlerTokenDetail)).Methods(http.MethodGet)

	apiV1.Handle("/register",
		HandlerFunc(HandlerUserRegister)).Methods(http.MethodPost)
	apiV1.Handle("/login",
		HandlerFunc(HandlerUserLogin)).Methods(http.MethodPost)

	return r

}
