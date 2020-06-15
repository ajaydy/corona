package routers

import (
	"corona/api"
	"corona/helpers"
	"database/sql"
)

var (
	dbPool              *sql.DB
	logger              *helpers.Logger
	coronaService       *api.CoronaModule
	countryService      *api.CountryModule
	continentService    *api.ContinentModule
	userService         *api.UserModule
	rateLimitService    *api.RateLimitModule
	subscriptionService *api.SubscriptionModule
	tokenService        *api.TokenModule
)

func Init(db *sql.DB, log *helpers.Logger) {
	dbPool = db
	logger = log
	coronaService = api.NewCoronaModule(dbPool, logger)
	countryService = api.NewCountryModule(dbPool, logger)
	continentService = api.NewContinentModule(dbPool, logger)
	userService = api.NewUserModule(dbPool, logger)
	rateLimitService = api.NewRateLimitModule(dbPool, logger)
	subscriptionService = api.NewSubscriptionModule(dbPool, logger)
	tokenService = api.NewTokenModule(dbPool, logger)
}
