package scrapper

//func TestGetCoronaData(t *testing.T) {
//
//	c := colly.NewCollector(
//		colly.AllowedDomains("www.worldometers.info"),
//	)
//
//	var datas []models.CoronaDataModel
//
//	c.OnRequest(func(r *colly.Request) {
//		log.Println("visiting", r.URL.String())
//	})
//
//	c.OnHTML("table#main_table_countries_today tbody", func(e *colly.HTMLElement) {
//
//		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
//
//			tds := el.ChildTexts("td")
//
//			for i := 0; i < len(tds); i++ {
//				if tds[i] == "" {
//					tds[i] = "0"
//				}
//			}
//
//			data := models.CoronaDataModel{
//				Name:           tds[1],
//				TotalCases:     tds[2],
//				NewCases:       tds[3],
//				TotalDeaths:    tds[4],
//				NewDeaths:      tds[5],
//				TotalRecovered: tds[6],
//				ActiveCases:    tds[8],
//				SeriousCases:   tds[9],
//				TotalTests:     tds[12],
//				Population:     tds[14],
//			}
//			if data.Name != "0" && data.Name != "Total:" {
//				datas = append(datas, data)
//			}
//
//		})
//
//	})
//
//	c.Visit("https://www.worldometers.info/coronavirus")
//
//	enc := json.NewEncoder(os.Stdout)
//	enc.SetIndent("", "  ")
//
//	enc.Encode(datas)
//
//}
//
//func TestGetDeaths(t *testing.T) {
//
//	c := colly.NewCollector(
//		colly.AllowedDomains("www.worldometers.info"),
//	)
//
//	c.OnHTML("div[id=maincounter-wrap]", func(e *colly.HTMLElement) {
//
//		e.ForEach(".maincounter-number", func(_ int, elem *colly.HTMLElement) {
//			fmt.Println(elem.Text)
//
//		})
//
//	})
//
//	c.Visit("https://www.worldometers.info/coronavirus")
//
//}
