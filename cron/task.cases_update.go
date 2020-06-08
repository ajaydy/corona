package cron

import (
	"context"
	"corona/api"
	"net/url"
)

type (
	TaskCasesUpdate struct {
		Name string
	}
)

const (
	cronCasesUpdate = "cases_update"
)

func (t *TaskCasesUpdate) Parse(opt url.Values) error {
	return nil
}

func (t TaskCasesUpdate) Run(ctx context.Context) {
	errLogger := logger.Err.WithField("cron", cronCasesUpdate)
	outLogger := logger.Out.WithField("cron", cronCasesUpdate)

	outLogger.Println("Updating data...")

	err := api.UpdateAllCorona(ctx)

	if err != nil {
		errLogger.WithError(err).Errorln("Error on updating cases.")
		return
	}
	outLogger.Println("Cases updated successfully")
}
