package models

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"time"
)

type (
	SubscriptionModel struct {
		Id               uuid.UUID
		SubscriptionType string
		RequestPerDay    int
		IsDelete         bool
		CreatedBy        uuid.UUID
		CreatedAt        time.Time
		UpdatedBy        uuid.NullUUID
		UpdatedAt        pq.NullTime
	}

	SubscriptionResponse struct {
		Id               uuid.UUID `json:"id"`
		SubscriptionType string    `json:"subscription_type"`
		RequestPerDay    int       `json:"request_per_day"`
		IsDelete         bool      `json:"is_delete"`
		CreatedBy        uuid.UUID `json:"created_by"`
		CreatedAt        time.Time `json:"created_at"`
		UpdatedBy        uuid.UUID `json:"updated_by"`
		UpdatedAt        time.Time `json:"updated_at"`
	}
)

func (s SubscriptionModel) Response() SubscriptionResponse {

	return SubscriptionResponse{
		Id:               s.Id,
		SubscriptionType: s.SubscriptionType,
		RequestPerDay:    s.RequestPerDay,
		IsDelete:         s.IsDelete,
		CreatedBy:        s.CreatedBy,
		CreatedAt:        s.CreatedAt,
		UpdatedBy:        s.UpdatedBy.UUID,
		UpdatedAt:        s.UpdatedAt.Time,
	}

}

func GetOneSubscription(ctx context.Context, db *sql.DB, id uuid.UUID) (SubscriptionModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			subscription_type,
			request_per_day,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			subscription
		WHERE 
			id = $1
	`)

	var subscription SubscriptionModel
	err := db.QueryRowContext(ctx, query, id).Scan(
		&subscription.Id,
		&subscription.SubscriptionType,
		&subscription.RequestPerDay,
		&subscription.IsDelete,
		&subscription.CreatedBy,
		&subscription.CreatedAt,
		&subscription.UpdatedBy,
		&subscription.UpdatedAt,
	)

	if err != nil {
		return SubscriptionModel{}, err
	}

	return subscription, nil

}

func GetAllSubscription(ctx context.Context, db *sql.DB) ([]SubscriptionModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			subscription_type,
			request_per_day,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM subscription
		`)

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var subscriptions []SubscriptionModel
	for rows.Next() {
		var subscription SubscriptionModel

		rows.Scan(
			&subscription.Id,
			&subscription.SubscriptionType,
			&subscription.RequestPerDay,
			&subscription.IsDelete,
			&subscription.CreatedBy,
			&subscription.CreatedAt,
			&subscription.UpdatedBy,
			&subscription.UpdatedAt,
		)

		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil

}

func (s *SubscriptionModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO subscription(
			subscription_type,
			request_per_day,
			created_by,
			created_at
		)VALUES(
			$1,$2,$3,now())
		RETURNING
			id, created_at
	`)

	err := db.QueryRowContext(ctx, query,
		s.SubscriptionType, s.RequestPerDay, s.CreatedBy).Scan(
		&s.Id, &s.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}
