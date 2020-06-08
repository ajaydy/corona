package routers

import (
	"corona/api"
	"corona/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerCountryList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCountryList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return countryService.List(ctx, filter)
}

func HandlerCountryDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	countryId, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCountryDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.CountryDetailParam{Id: countryId}

	return countryService.Detail(ctx, param)
}

func HandlerCountryAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.CountryAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerCountryAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return countryService.Add(ctx, param)
}
