package scrapper

import (
	"context"
	"corona/models"
	"github.com/gocolly/colly/v2"
	uuid "github.com/satori/go.uuid"
	"log"
	"strings"
)

func GetCoronaData(ctx context.Context) ([]models.CoronaModel, error) {

	c := colly.NewCollector(
		colly.AllowedDomains("www.worldometers.info"),
	)

	var datas []models.CoronaModel

	countries, err := models.GetAllCountry(ctx, dbPool)

	if err != nil {
		return nil, err
	}

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnHTML("table#main_table_countries_today tbody", func(e *colly.HTMLElement) {

		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {

			tds := el.ChildTexts("td")

			for i := 0; i < len(tds); i++ {
				if tds[i] == "" {
					tds[i] = "0"
				}
				if tds[i] == "Czechia" {
					tds[i] = "Czech Republic"
				}
				if tds[i] == "DRC" {
					tds[i] = "DR Congo"
				}
				if tds[i] == "Ivory Coast" {
					tds[i] = "Cote D'Ivoire"
				}
				if tds[i] == "CAR" {
					tds[i] = "Central African Republic"
				}
				if tds[i] == "UAE" {
					tds[i] = "United Arab Emirates"
				}
				if tds[i] == "UK" {
					tds[i] = "United Kingdom"
				}
				if tds[i] == "USA" {
					tds[i] = "United States"
				}
				if tds[i] == "S.Korea" {
					tds[i] = "South Korea"
				}
				if tds[i] == "USA" {
					tds[i] = "United States"
				}
			}

			var countryID uuid.UUID
			for _, country := range countries {

				if strings.EqualFold(country.Name, tds[1]) == true {
					countryID = country.Id
				}

			}

			data := models.CoronaModel{
				CountryId:      countryID,
				TotalCases:     tds[2],
				NewCases:       tds[3],
				TotalDeaths:    tds[4],
				NewDeaths:      tds[5],
				TotalRecovered: tds[6],
				ActiveCases:    tds[8],
				SeriousCases:   tds[9],
				TotalTests:     tds[12],
				Population:     tds[14],
			}

			if data.CountryId != uuid.Nil {
				datas = append(datas, data)
			}

		})

	})

	c.Visit("https://www.worldometers.info/coronavirus")

	return datas, nil

}
