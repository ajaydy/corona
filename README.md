# CoronaVirus Scrapper
This is a simple project to scrap corona-virus data from https://www.worldometers.info/coronavirus/. 

Data is being updated by using go-cron.



## Prerequisites
* Golang 1.12.17 or newer with Go Module Support

## Setup The Project

### Clone The Project
`git clone https://github.com/ajaydy/corona.git`

### Edit Config file
Copy ```example.config.toml``` to ```.config.toml```

`cp example.config.toml .config.toml`

Fill in your database, cron and app details.


### Run The Project
```go run main.go```