package middleware

import (
	"corona/helpers"
	"database/sql"
)

var (
	dbPool *sql.DB
	logger *helpers.Logger
)

func Init(db *sql.DB, log *helpers.Logger) {
	dbPool = db
	logger = log
}
