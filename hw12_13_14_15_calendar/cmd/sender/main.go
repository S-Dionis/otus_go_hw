package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/brocker/sender"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/spf13/viper"
)

var pathToConfig string

func init() {
	flag.StringVar(&pathToConfig, "config", "configs/sender_config.yaml", "Path to configuration file")
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
		fmt.Printf("Error reading config file, %s", err)
		os.Exit(1)
	}

	var rabbitConf sender.RabbitConf
	err = viper.Sub("rabbit").Unmarshal(&rabbitConf)
	if err != nil {
		fmt.Printf("Error unmarshalling config file, %s", err)
		os.Exit(1)
	}

	newSender := sender.NewSender(rabbitConf)
	err = newSender.Connect()
	if err != nil {
		return
	}

	messages, err := newSender.Consume(rabbitConf.QueueName)
	if err != nil {
		return
	}

	go func() {
		for d := range messages {
			message := fmt.Sprintf("Received a message: %s", d.Body)
			slog.Info(message)
		}
	}()
	select {}
}
