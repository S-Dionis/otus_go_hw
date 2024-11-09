package sender

import (
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
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

func GetChannel(s *Sender) *amqp.Channel {
	return s.channel
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

func (s *Sender) Consume(queueName string) (<-chan amqp.Delivery, error) {
	messages, err := s.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *Sender) Send(text string) {
	message := fmt.Sprintf("send notification: %s", text)
	slog.Info(message)
}
