package routers

import (
	"corona/api"
	"corona/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerSubscriptionList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	return subscriptionService.List(ctx)
}

func HandlerSubscriptionDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	subscriptionID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSubscriptionDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.SubscriptionDetailParam{Id: subscriptionID}

	return subscriptionService.Detail(ctx, param)
}
