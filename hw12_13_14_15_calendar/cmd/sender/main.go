package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/cmd/config"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/brocker/scheduler"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/brocker/sender"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/spf13/viper"
)

var pathToConfig string

func init() {
	fmt.Println("Init sender service")
	flag.StringVar(&pathToConfig, "config", "configs/sender_config.yaml", "Path to configuration file")
	flag.Parse()
}

func main() {
	err := logger.InitLogger("INFO")
	if err != nil {
		fmt.Println("init logger error:", err)
		return
	}

	slog.Info(fmt.Sprintf("set config file %s", pathToConfig))
	viper.SetConfigFile(pathToConfig)

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file, %s", err)
		os.Exit(1)
	}

	var rabbitConf sender.RabbitConf
	var database storage.Storage
	var db scheduler.DBType

	err = viper.Sub("rabbit").Unmarshal(&rabbitConf)
	if err != nil {
		fmt.Printf("Error unmarshalling config file, %s", err)
		os.Exit(1)
	}

	err = viper.Sub("db").Unmarshal(&db)
	if err != nil {
		slog.Error(fmt.Sprintf("Error unmarshalling config file, %s", err))
		os.Exit(1)
	}

	switch db.Type {
	case "memory":
		database = memorystorage.New()
	case "sql":
		var psqlConf config.PsqlConf
		err = viper.Sub("postgres").Unmarshal(&psqlConf)
		if err != nil {
			fmt.Printf("Error unmarshalling config file, %s", err)
			os.Exit(1)
		}

		database = sqlstorage.New(context.Background(), psqlConf)
	}

	slog.Info("Config read successfully, run sender service")

	newSender := sender.NewSender(&rabbitConf, &database)
	err = newSender.Connect()
	if err != nil {
		slog.Error(fmt.Sprintf("Connect to rabbitmq error: %e", err))
		return
	}

	slog.Info("Connect to rabbitmq successfully")
	messages, err := newSender.Consume(rabbitConf.QueueName)
	if err != nil {
		slog.Error(fmt.Sprintf("Receive message error: %e", err))
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	go func() {
		slog.Info("Sender is running ...")
		for d := range messages {
			err := newSender.Send(d.Body)
			if err != nil {
				slog.Error(fmt.Sprintf("Error sending message %e", err))
			}
		}
	}()
	select {}
}
