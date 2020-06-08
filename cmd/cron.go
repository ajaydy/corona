package cmd

import (
	"context"
	"corona/api"
	"corona/cron"
	"corona/scrapper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

var cronCmd = &cobra.Command{
	Use: "cron",
	PreRun: func(cmd *cobra.Command, args []string) {

		initDB()
		initLogger()
		scrapper.Init(dbPool, logger)

		cron.Init(logger, dbPool, cron.InitOption{
			CronTask: viper.GetStringMap("cron"),
		})

		api.Init(dbPool, logger)

	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		tasks, runningCron := cron.Run(ctx)
		logger.Out.Println("Cron up & running.")
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		logger.Out.Println("Stopping cron..")
		tasks.Clear()
		close(runningCron)
		logger.Out.Println("Cron stopped.")
	},
}
