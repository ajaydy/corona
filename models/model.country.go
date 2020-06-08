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
	CountryModel struct {
		Id          uuid.UUID
		ContinentId uuid.UUID
		Name        string
		Code        string
		CreatedBy   uuid.UUID
		CreatedAt   time.Time
		UpdatedBy   uuid.NullUUID
		UpdatedAt   pq.NullTime
	}

	CountryResponse struct {
		Id        uuid.UUID         `json:"id"`
		Continent ContinentResponse `json:"continent"`
		Name      string            `json:"name"`
		Code      string            `json:"code"`
		CreatedBy uuid.UUID         `json:"created_by"`
		CreatedAt time.Time         `json:"created_at"`
		UpdatedBy uuid.UUID         `json:"updated_by"`
		UpdatedAt time.Time         `json:"updated_at"`
	}
)

func (s CountryModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (CountryResponse, error) {

	continent, err := GetOneContinent(ctx, db, s.ContinentId)

	if err != nil {
		return CountryResponse{}, err
	}

	return CountryResponse{
		Id:        s.Id,
		Continent: continent.Response(),
		Name:      s.Name,
		Code:      s.Code,
		CreatedBy: s.CreatedBy,
		CreatedAt: s.CreatedAt,
		UpdatedBy: s.UpdatedBy.UUID,
		UpdatedAt: s.UpdatedAt.Time,
	}, nil
}

func GetOneCountry(ctx context.Context, db *sql.DB, id uuid.UUID) (CountryModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			continent_id,
			name,	
			code,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			country
		WHERE 
			id = $1
	`)

	var country CountryModel
	err := db.QueryRowContext(ctx, query, id).Scan(
		&country.Id,
		&country.ContinentId,
		&country.Name,
		&country.Code,
		&country.CreatedBy,
		&country.CreatedAt,
		&country.UpdatedBy,
		&country.UpdatedAt,
	)

	if err != nil {
		return CountryModel{}, err
	}

	return country, nil

}

func GetAllCountry(ctx context.Context, db *sql.DB) ([]CountryModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			continent_id,
			name,
			code,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			country
		`)

	rows, err := db.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var countries []CountryModel
	for rows.Next() {
		var country CountryModel

		rows.Scan(
			&country.Id,
			&country.ContinentId,
			&country.Name,
			&country.Code,
			&country.CreatedBy,
			&country.CreatedAt,
			&country.UpdatedBy,
			&country.UpdatedAt,
		)

		countries = append(countries, country)
	}

	return countries, nil
}

func GetAllCountries(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]CountryModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`WHERE LOWER(name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	query := fmt.Sprintf(`
		SELECT
			id,
			continent_id,
			name,
			code,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			country
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

	var countries []CountryModel
	for rows.Next() {
		var country CountryModel

		rows.Scan(
			&country.Id,
			&country.ContinentId,
			&country.Name,
			&country.Code,
			&country.CreatedBy,
			&country.CreatedAt,
			&country.UpdatedBy,
			&country.UpdatedAt,
		)

		countries = append(countries, country)
	}

	return countries, nil
}

func (s *CountryModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO country(
			continent_id,
			name,
			code,
			created_by,
			created_at
		)VALUES(
			$1,$2,$3,$4,now())
		RETURNING
			id, created_at
	`)

	err := db.QueryRowContext(ctx, query,
		s.ContinentId, s.Name, s.Code, s.CreatedBy).Scan(
		&s.Id, &s.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}
