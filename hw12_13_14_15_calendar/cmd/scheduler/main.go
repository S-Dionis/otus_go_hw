package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/brocker/scheduler"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/viper"
)

var pathToConfig string

func init() {
	flag.StringVar(&pathToConfig, "config", "configs/scheduler_config.yaml", "Path to configuration file")
}

func main() {
	err := logger.InitLogger("INFO")
	if err != nil {
		fmt.Println("init logger error:", err)
		return
	}

	viper.SetConfigFile(pathToConfig)

	err = viper.ReadInConfig()
	if err != nil {
		slog.Error(fmt.Sprintf("Error reading config file, %s", err))
		os.Exit(1)
	}

	var rabbitConf scheduler.RabbitConf
	var db scheduler.DBType

	err = viper.Sub("rabbit").Unmarshal(&rabbitConf)
	if err != nil {
		slog.Error(fmt.Sprintf("Error unmarshalling config file, %s", err))
		os.Exit(1)
	}

	err = viper.Sub("db").Unmarshal(&db)
	if err != nil {
		slog.Error(fmt.Sprintf("Error unmarshalling config file, %s", err))
		os.Exit(1)
	}
	var database storage.Storage

	switch db.Type {
	case "memory":
		database = memorystorage.New()
	case "sql":
		database = sqlstorage.New()
	}

	newScheduler := scheduler.NewScheduler(&rabbitConf, &database)
	slog.Info("Initializing scheduler...")
	err = newScheduler.Connect()
	if err != nil {
		slog.Error(fmt.Sprintf("Error connecting to RabbitMQ, %s", err))
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			slog.Info("Checking database for updates")

			err := newScheduler.DeleteItems()
			if err != nil {
				return
			}

			err = newScheduler.DatabaseMonitor(context.Background())
			if err != nil {
				return
			}
		}
	}()

	select {}
}
