/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"corona/api"
	"corona/helpers"
	"corona/middleware"
	"corona/routers"
	"corona/scrapper"
	"database/sql"
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	cfgFile string
	logger  *helpers.Logger
	dbPool  *sql.DB
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "corona",
	PreRun: func(cmd *cobra.Command, args []string) {
		initDB()
		initLogger()
		scrapper.Init(dbPool, logger)
		api.Init(dbPool, logger)
		helpers.Init(logger)
		routers.Init(dbPool, logger)
		middleware.Init(dbPool, logger)
	},
	Run: func(cmd *cobra.Command, args []string) {
		router := routers.InitHandlers()

		server := &http.Server{
			Addr:         fmt.Sprintf(":%d", viper.GetInt("app.port")),
			ReadTimeout:  time.Duration(viper.GetInt("app.read_timeout")) * time.Second,
			WriteTimeout: time.Duration(viper.GetInt("app.write_timeout")) * time.Second,
			Handler:      router,
		}

		idleConnsClosed := make(chan struct{})
		go func() {
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
			<-sigint
			timeout := time.Duration(viper.GetInt("app.shutdown_timeout")) * time.Second
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				logger.Out.WithError(err).Println("Server shutdown error.")
			}
			logger.Out.Println("Core server shutdown.")
			close(idleConnsClosed)
		}()

		logger.Out.Println(fmt.Sprintf(`Server Listen And Serve On Port : %d`, viper.GetInt("app.port")))
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Out.Println(fmt.Sprintf(`Error Listen And Serve : %v`, err))
			os.Exit(0)
		}
		<-idleConnsClosed

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(cronCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig, initDB)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .config.toml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("toml")

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".config" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigName(".config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initDB() {
	dbOptions := helpers.DBOptions{
		Host:        viper.GetString("database.host"),
		Port:        viper.GetInt("database.port"),
		Username:    viper.GetString("database.username"),
		Password:    viper.GetString("database.password"),
		DBName:      viper.GetString("database.name"),
		SSLCert:     viper.GetString("database.sslcert"),
		SSLKey:      viper.GetString("database.sslkey"),
		SSLRootCert: viper.GetString("database.sslrootcert"),
		SSLMode:     viper.GetString("database.sslmode"),
	}

	db, err := helpers.InitDB(dbOptions)

	if err != nil {
		logger.Err.Println(fmt.Sprintf("err connect : %v", err))
		os.Exit(0)
	}

	dbPool = db
}

func initLogger() {
	logger = helpers.NewLogger()
	logger.Out.Formatter = new(logrus.JSONFormatter)
	logger.Err.Formatter = new(logrus.JSONFormatter)
}
