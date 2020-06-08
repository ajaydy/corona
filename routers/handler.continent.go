package routers

import (
	"corona/api"
	"corona/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerContinentList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerContinentList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return continentService.List(ctx, filter)
}

func HandlerContinentDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	dataId, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerContinentDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.ContinentDetailParam{Id: dataId}

	return continentService.Detail(ctx, param)
}
