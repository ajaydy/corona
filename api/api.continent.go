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
	ContinentModule struct {
		db     *sql.DB
		logger *helpers.Logger
		name   string
	}

	ContinentDetailParam struct {
		Id uuid.UUID `json:"id"`
	}

	ContinentAddParam struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}
)

func NewContinentModule(db *sql.DB, logger *helpers.Logger) *ContinentModule {
	return &ContinentModule{
		db:     db,
		logger: logger,
		name:   "module/continent",
	}
}

func (s ContinentModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	continents, err := models.GetAllContinent(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllContinent",
			helpers.InternalServerError, http.StatusInternalServerError)
	}

	var continentResponses []models.ContinentResponse
	for _, continent := range continents {
		continentResponses = append(continentResponses, continent.Response())
	}

	return continentResponses, nil
}

func (s ContinentModule) Detail(ctx context.Context, param ContinentDetailParam) (interface{}, *helpers.Error) {

	continent, err := models.GetOneContinent(ctx, s.db, param.Id)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneContinent",
			helpers.InternalServerError, http.StatusInternalServerError)
	}

	return continent.Response(), nil

}

func (s ContinentModule) Add(ctx context.Context, param ContinentAddParam) (interface{}, *helpers.Error) {

	continent := models.ContinentModel{
		Name:      param.Name,
		Code:      param.Code,
		CreatedBy: uuid.NewV4(),
	}

	err := continent.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return continent.Response(), nil

}
