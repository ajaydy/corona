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
	CoronaModel struct {
		Id             uuid.UUID
		CountryId      uuid.UUID
		TotalCases     string
		NewCases       string
		TotalDeaths    string
		NewDeaths      string
		TotalRecovered string
		ActiveCases    string
		SeriousCases   string
		TotalTests     string
		Population     string
		CreatedBy      uuid.UUID
		CreatedAt      time.Time
		UpdatedBy      uuid.NullUUID
		UpdatedAt      pq.NullTime
	}

	CoronaResponse struct {
		Id             uuid.UUID       `json:"id"`
		Country        CountryResponse `json:"country"`
		TotalCases     string          `json:"total_cases"`
		NewCases       string          `json:"new_cases"`
		TotalDeaths    string          `json:"total_deaths"`
		NewDeaths      string          `json:"new_deaths"`
		TotalRecovered string          `json:"total_recovered"`
		ActiveCases    string          `json:"active_cases"`
		SeriousCases   string          `json:"serious_cases"`
		TotalTests     string          `json:"total_tests"`
		Population     string          `json:"population"`
		CreatedBy      uuid.UUID       `json:"created_by"`
		CreatedAt      time.Time       `json:"created_at"`
		UpdatedBy      uuid.UUID       `json:"updated_by"`
		UpdatedAt      time.Time       `json:"updated_at"`
	}
)

func (s CoronaModel) Response(ctx context.Context, db *sql.DB, logger *helpers.Logger) (CoronaResponse, error) {

	country, err := GetOneCountry(ctx, db, s.CountryId)

	if err != nil {
		logger.Err.Printf(`model.corona.go/GetOneCountry/%v`, err)
		return CoronaResponse{}, err
	}

	countryResponse, err := country.Response(ctx, db, logger)

	if err != nil {
		logger.Err.Printf(`model.corona.go/country.Response/%v`, err)
		return CoronaResponse{}, err
	}

	return CoronaResponse{
		Id:             s.Id,
		Country:        countryResponse,
		TotalCases:     s.TotalCases,
		NewCases:       s.NewCases,
		TotalDeaths:    s.TotalDeaths,
		NewDeaths:      s.NewDeaths,
		TotalRecovered: s.TotalRecovered,
		ActiveCases:    s.ActiveCases,
		SeriousCases:   s.SeriousCases,
		TotalTests:     s.TotalTests,
		Population:     s.Population,
		CreatedBy:      s.CreatedBy,
		CreatedAt:      s.CreatedAt,
		UpdatedBy:      s.UpdatedBy.UUID,
		UpdatedAt:      s.UpdatedAt.Time,
	}, nil
}

func GetOneCorona(ctx context.Context, db *sql.DB, id uuid.UUID) (CoronaModel, error) {

	query := fmt.Sprintf(`
		SELECT
			id,
			country_id,
			total_cases,
			new_cases,
			total_deaths,
			new_deaths,
			total_recovered,
			active_cases,
			serious_cases,
			total_tests,
			population,
			created_by,
			created_at,
			updated_by,
			updated_at
		FROM 
			corona_data
		WHERE 
			id = $1
	`)

	var data CoronaModel
	err := db.QueryRowContext(ctx, query, id).Scan(
		&data.Id,
		&data.CountryId,
		&data.TotalCases,
		&data.NewCases,
		&data.TotalDeaths,
		&data.NewDeaths,
		&data.TotalRecovered,
		&data.ActiveCases,
		&data.SeriousCases,
		&data.TotalTests,
		&data.Population,
		&data.CreatedBy,
		&data.CreatedAt,
		&data.UpdatedBy,
		&data.UpdatedAt,
	)

	if err != nil {
		return CoronaModel{}, err
	}

	return data, nil

}

func GetCoronaByCountry(ctx context.Context, db *sql.DB, country string) (CoronaModel, error) {

	query := fmt.Sprintf(`
		SELECT
			cd.id,
			country_id,
			total_cases,
			new_cases,
			total_deaths,
			new_deaths,
			total_recovered,
			active_cases,
			serious_cases,
			total_tests,
			population,
			cd.created_by,
			cd.created_at,
			cd.updated_by,
			cd.updated_at
		FROM 
			corona_data cd
		INNER JOIN
			country c
		ON 
			cd.country_id=c.id
		WHERE 
			LOWER(c.name) LIKE LOWER($1)
	`)

	var data CoronaModel
	err := db.QueryRowContext(ctx, query, country).Scan(
		&data.Id,
		&data.CountryId,
		&data.TotalCases,
		&data.NewCases,
		&data.TotalDeaths,
		&data.NewDeaths,
		&data.TotalRecovered,
		&data.ActiveCases,
		&data.SeriousCases,
		&data.TotalTests,
		&data.Population,
		&data.CreatedBy,
		&data.CreatedAt,
		&data.UpdatedBy,
		&data.UpdatedAt,
	)

	if err != nil {
		return CoronaModel{}, err
	}

	return data, nil

}

func GetAllCoronaByContinent(ctx context.Context, db *sql.DB, filter helpers.Filter, id uuid.UUID) (
	[]CoronaModel, error) {

	query := fmt.Sprintf(`
		SELECT
			cd.id,
			country_id,
			total_cases,
			new_cases,
			total_deaths,
			new_deaths,
			total_recovered,
			active_cases,
			serious_cases,
			total_tests,
			population,
			cd.created_by,
			cd.created_at,
			cd.updated_by,
			cd.updated_at
		FROM 
			corona_data cd
		INNER JOIN 
			country c
		ON
			cd.country_id = c.id
		WHERE 
			c.continent_id = $1
		ORDER BY 
			c.name %s
		LIMIT $2 OFFSET $3`, filter.Dir)

	rows, err := db.QueryContext(ctx, query, id, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var datas []CoronaModel
	for rows.Next() {
		var data CoronaModel

		rows.Scan(
			&data.Id,
			&data.CountryId,
			&data.TotalCases,
			&data.NewCases,
			&data.TotalDeaths,
			&data.NewDeaths,
			&data.TotalRecovered,
			&data.ActiveCases,
			&data.SeriousCases,
			&data.TotalTests,
			&data.Population,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)

		datas = append(datas, data)
	}

	return datas, nil

}

func GetAllCorona(ctx context.Context, db *sql.DB, filter helpers.Filter) ([]CoronaModel, error) {

	var searchQuery string

	if filter.Search != "" {
		searchQuery = fmt.Sprintf(`WHERE LOWER(c.name) LIKE LOWER('%%%s%%')`, filter.Search)
	}

	query := fmt.Sprintf(`
		SELECT
			cd.id,
			country_id,
			total_cases,
			new_cases,
			total_deaths,
			new_deaths,
			total_recovered,
			active_cases,
			serious_cases,
			total_tests,
			population,
			cd.created_by,
			cd.created_at,
			cd.updated_by,
			cd.updated_at
		FROM 
			corona_data cd
		INNER JOIN 
			country c
		ON
			cd.country_id = c.id
		%s
		ORDER BY 
			c.name %s
		LIMIT $1 OFFSET $2`, searchQuery, filter.Dir)

	rows, err := db.QueryContext(ctx, query, filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var datas []CoronaModel
	for rows.Next() {
		var data CoronaModel

		rows.Scan(
			&data.Id,
			&data.CountryId,
			&data.TotalCases,
			&data.NewCases,
			&data.TotalDeaths,
			&data.NewDeaths,
			&data.TotalRecovered,
			&data.ActiveCases,
			&data.SeriousCases,
			&data.TotalTests,
			&data.Population,
			&data.CreatedBy,
			&data.CreatedAt,
			&data.UpdatedBy,
			&data.UpdatedAt,
		)

		datas = append(datas, data)
	}

	return datas, nil

}

func (s *CoronaModel) Insert(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		INSERT INTO corona_data(
			country_id,
			total_cases,
			new_cases,
			total_deaths,
			new_deaths,
			total_recovered,
			active_cases,
			serious_cases,
			total_tests,
			population,
			created_by,
			created_at
		)VALUES(
			$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,now())
		RETURNING
			id, created_at
	`)

	err := db.QueryRowContext(ctx, query,
		s.CountryId, s.TotalCases, s.NewCases, s.TotalDeaths, s.NewDeaths, s.TotalRecovered, s.ActiveCases, s.SeriousCases,
		s.TotalTests, s.Population, s.CreatedBy).Scan(
		&s.Id, &s.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *CoronaModel) Update(ctx context.Context, db *sql.DB) error {

	query := fmt.Sprintf(`
		UPDATE corona_data
		SET
			total_cases=$1,
			new_cases=$2,
			total_deaths=$3,
			new_deaths=$4,
			total_recovered=$5,
			active_cases=$6,
			serious_cases=$7,
			total_tests=$8,
			population=$9,
			updated_at=NOW(),
			updated_by=$10
		WHERE 
			country_id=$11
		RETURNING 
			id,created_at,updated_at,created_by
	`)

	err := db.QueryRowContext(ctx, query,
		s.TotalCases, s.NewCases, s.TotalDeaths, s.NewDeaths, s.TotalRecovered, s.ActiveCases, s.SeriousCases, s.TotalTests,
		s.Population, s.UpdatedBy, s.CountryId).Scan(
		&s.Id, &s.CreatedAt, &s.UpdatedAt, &s.CreatedBy,
	)

	if err != nil {
		return err
	}

	return nil

}
