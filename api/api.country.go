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
	CountryModule struct {
		db     *sql.DB
		logger *helpers.Logger
		name   string
	}

	CountryDetailParam struct {
		Id uuid.UUID `json:"id"`
	}

	Country struct {
		ContinentName string `json:"continent_name"`
		Name          string `json:"name"`
		Code          string `json:"code"`
	}

	CountryAddParam struct {
		Country []Country `json:"country"`
	}
)

func NewCountryModule(db *sql.DB, logger *helpers.Logger) *CountryModule {
	return &CountryModule{
		db:     db,
		logger: logger,
		name:   "module/country",
	}
}

func (s CountryModule) Detail(ctx context.Context, param CountryDetailParam) (interface{}, *helpers.Error) {

	country, err := models.GetOneCountry(ctx, s.db, param.Id)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneCountry", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := country.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/country.response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}

func (s CountryModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	countries, err := models.GetAllCountries(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllCountry", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var countryResponses []models.CountryResponse

	for _, country := range countries {

		response, err := country.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/country.Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		countryResponses = append(countryResponses, response)

	}

	return countryResponses, nil
}

func (s CountryModule) Add(ctx context.Context, param CountryAddParam) (interface{}, *helpers.Error) {

	var countryResponses []models.CountryResponse

	for _, country := range param.Country {

		continent, err := models.GetOneContinentByName(ctx, s.db, country.ContinentName)

		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Add/GetOneContinentByName", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		country := models.CountryModel{
			ContinentId: continent.Id,
			Name:        country.Name,
			Code:        country.Code,
			CreatedBy:   uuid.NewV4(),
		}

		err = country.Insert(ctx, s.db)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		countryResponse, err := country.Response(ctx, s.db, s.logger)

		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Add/country.Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		countryResponses = append(countryResponses, countryResponse)

	}

	return countryResponses, nil

}
