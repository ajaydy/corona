# CoronaVirus Scrapper
* This is a simple project to scrap corona-virus data from https://www.worldometers.info/coronavirus/. 

* Data is being updated by using go-cron.



## Prerequisites
* Golang 1.12.17 or newer with Go Module Support

## Features

### Authentication 

* Using Token Based Authentication by generating a random token for user when they register and storing it in the database.

* Comparing header token with token in database to verify authentication.

### Rate Limit

* Using rate limit to limit the  amount of API calls for different users. For example, free users get 50 requests/day  and premium users get 100 requests/day .

## Routes


- /api/coronavirus: summary of all countries' cases .

- /api/coronavirus/[continent]: summary of all countries' cases in this [continent] .

- /api/coronavirus/[country]: a summary of [country] cases .

- /api/countries: all countries and their ISO codes .

- /api/continents: all continents and their codes .


## Setup The Project

### Clone The Project
`git clone https://github.com/ajaydy/corona.git`

### Edit Config file
Copy ```example.config.toml``` to ```.config.toml```

`cp example.config.toml .config.toml`

Fill in your database, cron and app details.


### Run The Project
```go run main.go```