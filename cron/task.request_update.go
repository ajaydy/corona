package cron

import (
	"context"
	"corona/api"
	"net/url"
)

type (
	TaskRequestUpdate struct {
		Name string
	}
)

const (
	cronRequestUpdate = "request_update"
)

func (t *TaskRequestUpdate) Parse(opt url.Values) error {
	return nil
}

func (t TaskRequestUpdate) Run(ctx context.Context) {
	errLogger := logger.Err.WithField("cron", cronRequestUpdate)
	outLogger := logger.Out.WithField("cron", cronRequestUpdate)

	outLogger.Println("Updating data...")

	err := api.UpdateRateLimit(ctx)

	if err != nil {
		errLogger.WithError(err).Errorln("Error on updating request.")
		return
	}
	outLogger.Println("Request updated successfully")
}
