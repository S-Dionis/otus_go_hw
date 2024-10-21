package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

type Sender struct {
	url     string
	channel *amqp.Channel
	conn    *amqp.Connection
}

type RabbitConf struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Pass         string `mapstructure:"password"`
	ExchangeName string `mapstructure:"exchange_name"`
	QueueName    string `mapstructure:"queue_name"`
	RoutingKey   string `mapstructure:"routing_key"`
}

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

	var rabbitConf RabbitConf
	err = viper.Sub("rabbit").Unmarshal(&rabbitConf)
	if err != nil {
		fmt.Printf("Error unmarshalling config file, %s", err)
		os.Exit(1)
	}

	sender := NewSender(rabbitConf)
	err = sender.Connect()
	if err != nil {
		return
	}

	messages, err := sender.channel.Consume(
		rabbitConf.QueueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
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

func NewSender(conf RabbitConf) *Sender {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s", conf.User, conf.Pass, conf.Host, conf.Port)
	return &Sender{
		url:     url,
		channel: nil,
		conn:    nil,
	}
}

func (s *Sender) Connect() error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	s.conn = conn
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	s.channel = ch

	return nil
}

func (s *Sender) Send(text string) {
	message := fmt.Sprintf("send notification: %s", text)
	slog.Info(message)
}
