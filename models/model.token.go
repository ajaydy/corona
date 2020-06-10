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
	TokenModel struct {
		Id        uuid.UUID
		UserId    uuid.UUID
		TokenKey  string
		ExpiredAt time.Time
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	TokenResponse struct {
		Id        uuid.UUID `json:"id"`
		UserId    uuid.UUID `json:"user_id"`
		TokenKey  string    `json:"token_key"`
		ExpiredAt time.Time `json:"expired_at"`
		CreatedBy uuid.UUID `json:"created_by"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedBy uuid.UUID `json:"updated_by"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func (t TokenModel) Response() TokenResponse {

	return TokenResponse{
		Id:        t.Id,
		UserId:    t.UserId,
		TokenKey:  t.TokenKey,
		ExpiredAt: t.ExpiredAt,
		CreatedBy: t.CreatedBy,
		CreatedAt: t.CreatedAt,
		UpdatedBy: t.UpdatedBy.UUID,
		UpdatedAt: t.UpdatedAt.Time,
	}

}

func GetOneToken(ctx context.Context, db *sql.DB, id uuid.UUID) (TokenModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			user_id
			token_key,
			expired_at,
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
