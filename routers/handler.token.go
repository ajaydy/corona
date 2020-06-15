package routers

import (
	"corona/api"
	"corona/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerTokenList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	return tokenService.List(ctx)
}

func HandlerTokenDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	tokenID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerTokenDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.TokenDetailParam{Id: tokenID}

	return tokenService.Detail(ctx, param)
}
