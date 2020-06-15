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
	TokenModel struct {
		Id        uuid.UUID
		UserId    uuid.UUID
		TokenKey  string
		ExpiredAt time.Time
		IsActive  bool
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	TokenResponse struct {
		Id        uuid.UUID    `json:"id"`
		User      UserResponse `json:"user"`
		TokenKey  string       `json:"token_key"`
		ExpiredAt time.Time    `json:"expired_at"`
		IsActive  bool         `json:"is_active"`
		CreatedBy uuid.UUID    `json:"created_by"`
		CreatedAt time.Time    `json:"created_at"`
		UpdatedBy uuid.UUID    `json:"updated_by"`
		UpdatedAt time.Time    `json:"updated_at"`
	}
)

func (t TokenModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (TokenResponse, error) {

	user, err := GetOneUser(ctx, db, t.UserId)

	if err != nil {
		logger.Err.Printf(`model.token.go/GetOneUser/%v`, err)
		return TokenResponse{}, err
	}

	userResponse, err := user.Response(ctx, db, logger)

	if err != nil {
		logger.Err.Printf(`model.token.go/user.Response/%v`, err)
		return TokenResponse{}, err
	}

	return TokenResponse{
		Id:        t.Id,
		User:      userResponse,
		TokenKey:  t.TokenKey,
		ExpiredAt: t.ExpiredAt,
		IsActive:  t.IsActive,
		CreatedBy: t.CreatedBy,
		CreatedAt: t.CreatedAt,
		UpdatedBy: t.UpdatedBy.UUID,
		UpdatedAt: t.UpdatedAt.Time,
	}, nil

}

func GetOneToken(ctx context.Context, db *sql.DB, id uuid.UUID) (TokenModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			user_id,
			token_key,
			expired_at,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			token
		WHERE 
			id = $1
	`)

	var token TokenModel
	err := db.QueryRowContext(ctx, query, id).Scan(
		&token.Id,
		&token.UserId,
		&token.TokenKey,
		&token.ExpiredAt,
		&token.IsActive,
		&token.CreatedBy,
		&token.CreatedAt,
		&token.UpdatedBy,
		&token.UpdatedAt,
	)

	if err != nil {
		return TokenModel{}, err
	}

	return token, nil

}

func GetOneTokenByTokenKey(ctx context.Context, db *sql.DB, tokenKey string) (TokenModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			user_id,
			token_key,
			expired_at,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			token
		WHERE 
			token_key = $1
	`)

	var token TokenModel
	err := db.QueryRowContext(ctx, query, tokenKey).Scan(
		&token.Id,
		&token.UserId,
		&token.TokenKey,
		&token.ExpiredAt,
		&token.IsActive,
		&token.CreatedBy,
		&token.CreatedAt,
		&token.UpdatedBy,
		&token.UpdatedAt,
	)

	if err != nil {
		return TokenModel{}, err
	}

	return token, nil

}

func GetAllToken(ctx context.Context, db *sql.DB) ([]TokenModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			user_id,
			token_key,
			expired_at,
			is_active,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			token
		`)

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tokens []TokenModel
	for rows.Next() {
		var token TokenModel

		rows.Scan(
			&token.Id,
			&token.UserId,
			&token.TokenKey,
			&token.ExpiredAt,
			&token.IsActive,
			&token.CreatedBy,
			&token.CreatedAt,
			&token.UpdatedBy,
			&token.UpdatedAt,
		)

		tokens = append(tokens, token)
	}

	return tokens, nil

}

func (t *TokenModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO token(
			user_id,
			token_key,
			expired_at,
			created_by,
			created_at
		)VALUES(
			$1,$2,$3,$4,now())
		RETURNING
			id, created_at
	`)

	err := db.QueryRowContext(ctx, query,
		t.UserId, t.TokenKey, t.ExpiredAt, t.CreatedBy).Scan(
		&t.Id, &t.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}
