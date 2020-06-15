package models

import (
	"context"
	"corona/helpers"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
	"time"
)

type (
	UserModel struct {
		Id             uuid.UUID
		SubscriptionId uuid.UUID
		Name           string
		Email          string
		Password       string
		IsActive       bool
		CreatedBy      uuid.UUID
		CreatedAt      time.Time
		UpdatedBy      uuid.NullUUID
		UpdatedAt      pq.NullTime
	}

	UserResponse struct {
		Id           uuid.UUID            `json:"id"`
		Subscription SubscriptionResponse `json:"subscription"`
		Name         string               `json:"name"`
		Email        string               `json:"email"`
		IsActive     bool                 `json:"is_active"`
		CreatedBy    uuid.UUID            `json:"created_by"`
		CreatedAt    time.Time            `json:"created_at"`
		UpdatedBy    uuid.UUID            `json:"updated_by"`
		UpdatedAt    time.Time            `json:"updated_at"`
	}
)

func (u UserModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (UserResponse, error) {

	subscription, err := GetOneSubscription(ctx, db, u.SubscriptionId)

	if err != nil {
		logger.Err.Printf(`model.user.go/GetOneSubscription/%v`, err)
		return UserResponse{}, err
	}

	return UserResponse{
		Id:           u.Id,
		Subscription: subscription.Response(),
		Name:         u.Name,
		Email:        u.Email,
		IsActive:     u.IsActive,
		CreatedBy:    u.CreatedBy,
		CreatedAt:    u.CreatedAt,
		UpdatedBy:    u.UpdatedBy.UUID,
		UpdatedAt:    u.UpdatedAt.Time,
	}, nil

}

func GetOneUser(ctx context.Context, db *sql.DB, id uuid.UUID) (UserModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			subscription_id,
			name,
			email,
			password,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			"user"
		WHERE 
			id = $1
	`)

	var user UserModel
	err := db.QueryRowContext(ctx, query, id).Scan(
		&user.Id,
		&user.SubscriptionId,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.CreatedBy,
		&user.CreatedAt,
		&user.UpdatedBy,
		&user.UpdatedAt,
	)

	if err != nil {
		return UserModel{}, err
	}

	return user, nil

}

func GetOneUserByEmail(ctx context.Context, db *sql.DB, email string) (UserModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			subscription_id,
			name,
			email,
			password,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			"user"
		WHERE 
			email=$1
	`)

	var user UserModel
	err := db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.SubscriptionId,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsActive,
		&user.CreatedBy,
		&user.CreatedAt,
		&user.UpdatedBy,
		&user.UpdatedAt,
	)

	if err != nil {
		return UserModel{}, err
	}

	return user, nil

}

func GetAllUser(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]UserModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`WHERE LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			subscription_id,
			name,
			email,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM user
		%s 
		ORDER BY
			name  %s
		LIMIT $1 OFFSET $2`,
		searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []UserModel
	for rows.Next() {
		var user UserModel

		rows.Scan(
			&user.Id,
			&user.SubscriptionId,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.IsActive,
			&user.CreatedBy,
			&user.CreatedAt,
			&user.UpdatedBy,
			&user.UpdatedAt,
		)

		users = append(users, user)
	}

	return users, nil

}

func (u *UserModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO "user"(
			subscription_id,
			name,
			email,
			password,
			created_by,
			created_at
		)VALUES(
			$1,$2,$3,$4,$5,now())
		RETURNING
			id, created_at,is_active
	`)

	err := db.QueryRowContext(ctx, query,
		u.SubscriptionId, u.Name, u.Email, u.Password, u.CreatedBy).Scan(
		&u.Id, &u.CreatedAt, &u.IsActive,
	)

	if err != nil {
		return err
	}

	return nil

}
