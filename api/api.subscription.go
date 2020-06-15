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
	SubscriptionModule struct {
		db     *sql.DB
		logger *helpers.Logger
		name   string
	}

	SubscriptionDetailParam struct {
		Id uuid.UUID `json:"id"`
	}

	SubscriptionAddParam struct {
		SubscriptionType string `json:"subscription_type"`
		RequestPerDay    int    `json:"request_per_day"`
	}
)

func NewSubscriptionModule(db *sql.DB, logger *helpers.Logger) *SubscriptionModule {
	return &SubscriptionModule{
		db:     db,
		logger: logger,
		name:   "module/subscription",
	}
}

func (s SubscriptionModule) List(ctx context.Context) (interface{}, *helpers.Error) {

	subscriptions, err := models.GetAllSubscription(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllSubscription",
			helpers.InternalServerError, http.StatusInternalServerError)
	}

	var subscriptionResponses []models.SubscriptionResponse
	for _, subscription := range subscriptions {
		subscriptionResponses = append(subscriptionResponses, subscription.Response())
	}

	return subscriptionResponses, nil
}

func (s SubscriptionModule) Detail(ctx context.Context, param SubscriptionDetailParam) (interface{}, *helpers.Error) {

	subscription, err := models.GetOneSubscription(ctx, s.db, param.Id)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneSubscription",
			helpers.InternalServerError, http.StatusInternalServerError)
	}

	return subscription.Response(), nil

}

func (s SubscriptionModule) Add(ctx context.Context, param SubscriptionAddParam) (interface{}, *helpers.Error) {

	subscription := models.SubscriptionModel{
		SubscriptionType: param.SubscriptionType,
		RequestPerDay:    param.RequestPerDay,
		CreatedBy:        uuid.NewV4(),
	}

	err := subscription.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return subscription.Response(), nil

}
