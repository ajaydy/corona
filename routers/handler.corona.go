package routers

import (
	"corona/api"
	"corona/helpers"
	"github.com/gorilla/mux"
	"net/http"
)

func HandlerCoronaList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCoronaList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return coronaService.List(ctx, filter)
}

func HandlerCoronaAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	return coronaService.Add(ctx)
}

func HandlerCoronaByCountry(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	country := params["country"]

	param := api.ByCountryParam{Country: country}

	return coronaService.ByCountry(ctx, param)
}

func HandlerCoronaByContinent(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCoronaByContinent/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	continent := params["continent"]

	param := api.ByContinentParam{Continent: continent}

	return coronaService.ByContinent(ctx, filter, param)
}
