package api

import (
	"context"
	"corona/helpers"
	"corona/models"
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"time"
)

type (
	TokenModule struct {
		db     *sql.DB
		logger *helpers.Logger
		name   string
	}

	TokenDetailParam struct {
		Id uuid.UUID `json:"id"`
	}

	TokenAddParam struct {
		UserId    uuid.UUID `json:"user_id"`
		TokenKey  string    `json:"token_key"`
		ExpiredAt time.Time `json:"expired_at"`
	}
)

func NewTokenModule(db *sql.DB, logger *helpers.Logger) *TokenModule {
	return &TokenModule{
		db:     db,
		logger: logger,
		name:   "module/token",
	}
}

func (t TokenModule) List(ctx context.Context) (interface{}, *helpers.Error) {

	tokens, err := models.GetAllToken(ctx, t.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, t.name, "List/GetAllToken",
			helpers.InternalServerError, http.StatusInternalServerError)
	}

	var tokenResponses []models.TokenResponse
	for _, token := range tokens {
		response, err := token.Response(ctx, t.db, t.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, t.name, "List/token.Response",
				helpers.InternalServerError, http.StatusInternalServerError)
		}
		tokenResponses = append(tokenResponses, response)
	}

	return tokenResponses, nil
}

func (t TokenModule) Detail(ctx context.Context, param TokenDetailParam) (interface{}, *helpers.Error) {

	token, err := models.GetOneToken(ctx, t.db, param.Id)

	if err != nil {
		return nil, helpers.ErrorWrap(err, t.name, "Detail/GetOneToken",
			helpers.InternalServerError, http.StatusInternalServerError)
	}

	response, err := token.Response(ctx, t.db, t.logger)
	if err != nil {
		return nil, helpers.ErrorWrap(err, t.name, "Detail/token.Response",
			helpers.InternalServerError, http.StatusInternalServerError)
	}

	return response, nil

}

func (t TokenModule) Add(ctx context.Context, param TokenAddParam) (interface{}, *helpers.Error) {

	token := models.TokenModel{
		UserId:    param.UserId,
		TokenKey:  param.TokenKey,
		ExpiredAt: param.ExpiredAt,
		CreatedBy: uuid.NewV4(),
	}

	err := token.Insert(ctx, t.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, t.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := token.Response(ctx, t.db, t.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, t.name, "Add/token.Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}
