package cron

import (
	"context"
	"corona/helpers"
	"database/sql"
	"github.com/jasonlvhit/gocron"
	"net/url"
	"strconv"
)

type (
	InitOption struct {
		CronTask map[string]interface{}
	}

	TaskRunner interface {
		Parse(opt url.Values) error
		Run(ctx context.Context)
	}
)

var mapTasks = map[string]TaskRunner{
	cronCasesUpdate:   &TaskCasesUpdate{},
	cronRequestUpdate: &TaskRequestUpdate{},
}

var (
	logger *helpers.Logger
	cfg    InitOption
	dbPool *sql.DB
)

func Init(lg *helpers.Logger, db *sql.DB, opt InitOption) {
	logger = lg
	cfg = opt
	dbPool = db
}

func Run(ctx context.Context) (*gocron.Scheduler, chan bool) {

	s := gocron.NewScheduler()

	runningCron := make(chan bool)
	taskCfg := cfg.CronTask["task"].([]interface{})

	for _, v := range taskCfg {

		name := v.(map[string]interface{})["name"].(string)
		interval, _ := strconv.Atoi(v.(map[string]interface{})["interval"].(string))
		unit := v.(map[string]interface{})["unit"].(string)
		time := v.(map[string]interface{})["time"].(string)

		switch unit {
		case "day":
			s.Every(uint64(interval)).Day().At(time).Do(mapTasks[name].Run, ctx)
		case "seconds":
			s.Every(uint64(interval)).Seconds().At(time).Do(mapTasks[name].Run, ctx)
		default:
			logger.Err.WithField("name", name).Println("Missing cron unit.")
		}
	}

	runningCron = s.Start()

	return s, runningCron
}
