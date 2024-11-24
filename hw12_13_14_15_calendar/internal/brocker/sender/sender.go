package sender

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/entities"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Sender struct {
	url     string
	channel *amqp.Channel
	conn    *amqp.Connection
	conf    *RabbitConf
	storage *storage.Storage
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

func NewSender(conf *RabbitConf, storage *storage.Storage) *Sender {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s", conf.User, conf.Pass, conf.Host, conf.Port)
	return &Sender{
		url:     url,
		channel: nil,
		conn:    nil,
		conf:    conf,
		storage: storage,
	}
}

func (s *Sender) Connect() error {
	conn, err := amqp.Dial(s.url)
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
	err := s.channel.ExchangeDeclare(
		s.conf.ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	_, err = s.channel.QueueDeclare(
		s.conf.QueueName, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return nil, err
	}

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

func (s *Sender) Send(body []byte) error {
	var event entities.Event
	err := json.Unmarshal(body, &event)
	if err != nil {
		slog.Error(fmt.Sprintf("Error unmarshalling %v event: %s", body, err))
		return err
	}

	message := fmt.Sprintf("send notification: %v", event)
	slog.Info(message)
	event.Notified = true
	storage := *s.storage
	err = storage.Change(&event)
	if err != nil {
		slog.Error(fmt.Sprintf("Error update event %v: %s", body, err))
		return err
	}

	return nil
}
