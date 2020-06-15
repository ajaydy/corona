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
	RateLimitModel struct {
		Id           uuid.UUID
		UserId       uuid.UUID
		TotalRequest int
		IsDelete     bool
		CreatedBy    uuid.UUID
		CreatedAt    time.Time
		UpdatedBy    uuid.NullUUID
		UpdatedAt    pq.NullTime
	}

	RateLimitResponse struct {
		Id           uuid.UUID    `json:"id"`
		User         UserResponse `json:"user"`
		TotalRequest int          `json:"total_request"`
		IsDelete     bool         `json:"is_delete"`
		CreatedBy    uuid.UUID    `json:"created_by"`
		CreatedAt    time.Time    `json:"created_at"`
		UpdatedBy    uuid.UUID    `json:"updated_by"`
		UpdatedAt    time.Time    `json:"updated_at"`
	}
)

func (r RateLimitModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (RateLimitResponse, error) {

	user, err := GetOneUser(ctx, db, r.UserId)
	if err != nil {
		logger.Err.Printf(`model.rate.limit.go/GetOneUser/%v`, err)
		return RateLimitResponse{}, err
	}

	userResponse, err := user.Response(ctx, db, logger)
	if err != nil {
		logger.Err.Printf(`model.rate.limit.go/user.Response/%v`, err)
		return RateLimitResponse{}, err
	}

	return RateLimitResponse{
		Id:           r.Id,
		User:         userResponse,
		TotalRequest: r.TotalRequest,
		IsDelete:     r.IsDelete,
		CreatedBy:    r.CreatedBy,
		CreatedAt:    r.CreatedAt,
		UpdatedBy:    r.UpdatedBy.UUID,
		UpdatedAt:    r.UpdatedAt.Time,
	}, nil

}

func GetOneRateLimit(ctx context.Context, db *sql.DB, id uuid.UUID) (RateLimitModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			user_id,
			total_request,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			rate_limit
		WHERE 
			id = $1
	`)

	var limit RateLimitModel
	err := db.QueryRowContext(ctx, query, id).Scan(
		&limit.Id,
		&limit.UserId,
		&limit.TotalRequest,
		&limit.IsDelete,
		&limit.CreatedBy,
		&limit.CreatedAt,
		&limit.UpdatedBy,
		&limit.UpdatedAt,
	)

	if err != nil {
		return RateLimitModel{}, err
	}

	return limit, nil

}

func GetOneRateLimitByUserId(ctx context.Context, db *sql.DB, userId uuid.UUID) (RateLimitModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			user_id,
			total_request,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			rate_limit
		WHERE 
			user_id = $1
	`)

	var limit RateLimitModel
	err := db.QueryRowContext(ctx, query, userId).Scan(
		&limit.Id,
		&limit.UserId,
		&limit.TotalRequest,
		&limit.IsDelete,
		&limit.CreatedBy,
		&limit.CreatedAt,
		&limit.UpdatedBy,
		&limit.UpdatedAt,
	)

	if err != nil {
		return RateLimitModel{}, err
	}

	return limit, nil

}

func GetAllRateLimit(ctx context.Context, db *sql.DB) ([]RateLimitModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			user_id,
			total_request,
			is_delete,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			rate_limit
	`)

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var limits []RateLimitModel
	for rows.Next() {
		var limit RateLimitModel

		rows.Scan(
			&limit.Id,
			&limit.UserId,
			&limit.TotalRequest,
			&limit.IsDelete,
			&limit.CreatedBy,
			&limit.CreatedAt,
			&limit.UpdatedBy,
			&limit.UpdatedAt,
		)

		limits = append(limits, limit)
	}

	return limits, nil
}

func (r *RateLimitModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO rate_limit(
			user_id,
			total_request,
			created_by,
			created_at
		)VALUES(
			$1,$2,$3,now())
		RETURNING
			id, created_at
	`)

	err := db.QueryRowContext(ctx, query,
		r.UserId, r.TotalRequest, r.CreatedBy).Scan(
		&r.Id, &r.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}

func (r *RateLimitModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE rate_limit
		SET
			user_id = $1,
			total_request=$2,
			updated_at=NOW(),
			updated_by=$3
		WHERE 
			id=$4
		RETURNING
			id,created_at,updated_at,created_by
	`)

	err := db.QueryRowContext(ctx, query, r.UserId, r.TotalRequest, r.UpdatedBy, r.Id).Scan(
		&r.Id, &r.CreatedAt, &r.UpdatedAt, &r.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil

}
