package routers

import (
	"corona/api"
	"corona/helpers"
	"database/sql"
)

var (
	dbPool           *sql.DB
	logger           *helpers.Logger
	coronaService    *api.CoronaModule
	countryService   *api.CountryModule
	continentService *api.ContinentModule
)

func Init(db *sql.DB, log *helpers.Logger) {
	dbPool = db
	logger = log
	coronaService = api.NewCoronaModule(dbPool, logger)
	countryService = api.NewCountryModule(dbPool, logger)
	continentService = api.NewContinentModule(dbPool, logger)
}
