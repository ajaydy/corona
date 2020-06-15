package routers

import (
	"corona/api"
	"corona/helpers"
	"net/http"
)

func HandlerUserRegister(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.UserRegisterParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerUserRegister/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return userService.Register(ctx, param)
}

func HandlerUserLogin(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.UserLoginParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerUserLogin/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	return userService.Login(ctx, param)
}
