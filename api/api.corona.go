package api

import (
	"context"
	"corona/helpers"
	"corona/models"
	"corona/scrapper"
	"database/sql"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type (
	CoronaModule struct {
		db     *sql.DB
		logger *helpers.Logger
		name   string
	}

	CoronaDetailParam struct {
		Id uuid.UUID `json:"id"`
	}

	ByCountryParam struct {
		Country string `json:"country"`
	}

	ByContinentParam struct {
		Continent string `json:"continent"`
	}
)

func NewCoronaModule(db *sql.DB, logger *helpers.Logger) *CoronaModule {
	return &CoronaModule{
		db:     db,
		logger: logger,
		name:   "module/corona",
	}
}

func (s CoronaModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	datas, err := models.GetAllCorona(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllCoronaData", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var dataResponses []models.CoronaResponse
	for _, data := range datas {

		response, err := data.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		dataResponses = append(dataResponses, response)

	}

	return dataResponses, nil
}

func (s CoronaModule) Add(ctx context.Context) (interface{}, *helpers.Error) {

	datas, err := scrapper.GetCoronaData(ctx)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GetCoronaData", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var dataResponses []models.CoronaResponse

	for _, data := range datas {

		data := models.CoronaModel{
			CountryId:      data.CountryId,
			TotalCases:     data.TotalCases,
			NewCases:       data.NewCases,
			TotalDeaths:    data.TotalDeaths,
			NewDeaths:      data.NewDeaths,
			TotalRecovered: data.TotalRecovered,
			ActiveCases:    data.ActiveCases,
			SeriousCases:   data.SeriousCases,
			TotalTests:     data.TotalTests,
			Population:     data.Population,
			CreatedBy:      uuid.NewV4(),
		}

		err := data.Insert(ctx, s.db)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		response, err := data.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		dataResponses = append(dataResponses, response)
	}

	return dataResponses, nil

}

func UpdateAllCorona(ctx context.Context) error {

	datas, err := scrapper.GetCoronaData(ctx)

	if err != nil {
		return err
	}

	for _, data := range datas {

		data := models.CoronaModel{
			CountryId:      data.CountryId,
			TotalCases:     data.TotalCases,
			NewCases:       data.NewCases,
			TotalDeaths:    data.TotalDeaths,
			NewDeaths:      data.NewDeaths,
			TotalRecovered: data.TotalRecovered,
			ActiveCases:    data.ActiveCases,
			SeriousCases:   data.SeriousCases,
			TotalTests:     data.TotalTests,
			Population:     data.Population,
			UpdatedBy: uuid.NullUUID{
				UUID:  uuid.NewV4(),
				Valid: true,
			},
		}

		err := data.Update(ctx, dbPool)

		if err != nil {
			return err
		}

	}

	return nil
}

func (s CoronaModule) ByCountry(ctx context.Context, param ByCountryParam) (interface{}, *helpers.Error) {

	data, err := models.GetCoronaByCountry(ctx, s.db, param.Country)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "ByCountry/GetCoronaDataByCountry",
			helpers.InternalServerError, http.StatusInternalServerError)
	}

	response, err := data.Response(ctx, s.db, s.logger)

	return response, nil

}

func (s CoronaModule) ByContinent(ctx context.Context, filter helpers.Filter, param ByContinentParam) (
	interface{}, *helpers.Error) {

	continent, err := models.GetOneContinentByName(ctx, s.db, param.Continent)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "ByContinent/GetOneContinentByName",
			helpers.InternalServerError, http.StatusInternalServerError)
	}

	datas, err := models.GetAllCoronaByContinent(ctx, s.db, filter, continent.Id)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "ByContinent/GetAllCoronaByContinent",
			helpers.InternalServerError, http.StatusInternalServerError)
	}

	var responses []models.CoronaResponse

	for _, data := range datas {

		response, err := data.Response(ctx, s.db, s.logger)

		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "ByContinent/Response",
				helpers.InternalServerError, http.StatusInternalServerError)
		}

		responses = append(responses, response)

	}

	return responses, nil

}
