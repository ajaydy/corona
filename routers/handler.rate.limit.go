package routers

import (
	"corona/api"
	"corona/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerRateLimitList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	return rateLimitService.List(ctx)
}

func HandlerRateLimitDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	limitID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerRateLimitDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.RateLimitDetailParam{Id: limitID}

	return rateLimitService.Detail(ctx, param)
}
