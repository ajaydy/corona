package api

import (
	"context"
	"corona/helpers"
	"corona/models"
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type (
	RateLimitModule struct {
		db     *sql.DB
		logger *helpers.Logger
		name   string
	}

	RateLimitDetailParam struct {
		Id uuid.UUID `json:"id"`
	}

	RateLimitAddParam struct {
		UserId       uuid.UUID `json:"user_id"`
		TotalRequest int       `json:"total_request"`
	}
)

func NewRateLimitModule(db *sql.DB, logger *helpers.Logger) *RateLimitModule {
	return &RateLimitModule{
		db:     db,
		logger: logger,
		name:   "module/rate.limit",
	}
}

func (r RateLimitModule) Detail(ctx context.Context, param RateLimitDetailParam) (interface{}, *helpers.Error) {

	limit, err := models.GetOneRateLimit(ctx, r.db, param.Id)

	if err != nil {
		return nil, helpers.ErrorWrap(err, r.name, "Detail/GetOneRateLimit", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := limit.Response(ctx, r.db, r.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, r.name, "Detail/limit.Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}

func (r RateLimitModule) List(ctx context.Context) (interface{}, *helpers.Error) {

	limits, err := models.GetAllRateLimit(ctx, r.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, r.name, "List/GetAllRateLimit", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var limitResponse []models.RateLimitResponse

	for _, limit := range limits {

		response, err := limit.Response(ctx, r.db, r.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, r.name, "List/limit.Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		limitResponse = append(limitResponse, response)

	}

	return limitResponse, nil
}

func (r RateLimitModule) Add(ctx context.Context, param RateLimitAddParam) (interface{}, *helpers.Error) {

	limit := models.RateLimitModel{
		UserId:       param.UserId,
		TotalRequest: param.TotalRequest,
		CreatedBy:    uuid.NewV4(),
	}

	err := limit.Insert(ctx, r.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, r.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := limit.Response(ctx, r.db, r.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, r.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}

func UpdateRateLimit(ctx context.Context) error {

	limits, err := models.GetAllRateLimit(ctx, dbPool)

	if err != nil {
		return err
	}

	for _, limit := range limits {

		limit.TotalRequest = 0

		limit.UpdatedBy = uuid.NullUUID{
			UUID:  uuid.NewV4(),
			Valid: true,
		}

		err := limit.Update(ctx, dbPool)

		if err != nil {
			return err
		}

	}

	return nil

}
