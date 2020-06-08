package routers

import (
	"corona/helpers"
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

	apiV1.Handle("/coronavirus",
		HandlerFunc(HandlerCoronaAdd)).Methods(http.MethodPost)
	apiV1.Handle("/coronavirus",
		HandlerFunc(HandlerCoronaList)).Methods(http.MethodGet)
	apiV1.Handle("/coronavirus/continents/{continent}",
		HandlerFunc(HandlerCoronaByContinent)).Methods(http.MethodGet)
	apiV1.Handle("/coronavirus/countries/{country}",
		HandlerFunc(HandlerCoronaByCountry)).Methods(http.MethodGet)

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

	return r

}
