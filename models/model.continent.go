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
	ContinentModel struct {
		Id        uuid.UUID
		Name      string
		Code      string
		CreatedBy uuid.UUID
		CreatedAt time.Time
		UpdatedBy uuid.NullUUID
		UpdatedAt pq.NullTime
	}

	ContinentResponse struct {
		Id        uuid.UUID `json:"id"`
		Name      string    `json:"name"`
		Code      string    `json:"code"`
		CreatedBy uuid.UUID `json:"created_by"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedBy uuid.UUID `json:"updated_by"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func (s ContinentModel) Response() ContinentResponse {

	return ContinentResponse{
		Id:        s.Id,
		Name:      s.Name,
		Code:      s.Code,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}

}

func GetOneContinent(ctx context.Context, db *sql.DB, id uuid.UUID) (ContinentModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,	
			code,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			continent
		WHERE 
			id = $1
	`)

	var continent ContinentModel
	err := db.QueryRowContext(ctx, query, id).Scan(
		&continent.Id,
		&continent.Name,
		&continent.Code,
		&continent.CreatedBy,
		&continent.CreatedAt,
		&continent.UpdatedBy,
		&continent.UpdatedAt,
	)

	if err != nil {
		return ContinentModel{}, err
	}

	return continent, nil

}

func GetOneContinentByName(ctx context.Context, db *sql.DB, name string) (ContinentModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			name,	
			code,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			continent
		WHERE 
			LOWER(name) = LOWER($1)
	`)

	var continent ContinentModel
	err := db.QueryRowContext(ctx, query, name).Scan(
		&continent.Id,
		&continent.Name,
		&continent.Code,
		&continent.CreatedBy,
		&continent.CreatedAt,
		&continent.UpdatedBy,
		&continent.UpdatedAt,
	)

	if err != nil {
		return ContinentModel{}, err
	}

	return continent, nil

}

func GetAllContinent(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]ContinentModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`WHERE LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			name,
			code,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM continent
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

	var continents []ContinentModel
	for rows.Next() {
		var continent ContinentModel

		rows.Scan(
			&continent.Id,
			&continent.Name,
			&continent.Code,
			&continent.CreatedBy,
			&continent.CreatedAt,
			&continent.UpdatedBy,
			&continent.UpdatedAt,
		)

		continents = append(continents, continent)
	}

	return continents, nil

}

func (s *ContinentModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO continent(
			name,
			code,
			created_by,
			created_at
		)VALUES(
			$1,$2,$3,now())
		RETURNING
			id, created_at
	`)

	err := db.QueryRowContext(ctx, query,
		s.Name, s.Code, s.CreatedBy).Scan(
		&s.Id, &s.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}
