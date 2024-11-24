package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/entities"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitConf struct {
	Host         string `mapstructure:"host"`
	Port         string `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Pass         string `mapstructure:"password"`
	ExchangeName string `mapstructure:"exchange_name"`
	QueueName    string `mapstructure:"queue_name"`
	RoutingKey   string `mapstructure:"routing_key"`
}

type DBType struct {
	Type string `mapstructure:"type"`
}

type Scheduler struct {
	conf    *RabbitConf
	storage *storage.Storage
	channel *amqp.Channel
	conn    *amqp.Connection
}

func NewScheduler(conf *RabbitConf, storage *storage.Storage) *Scheduler {
	return &Scheduler{
		conf:    conf,
		storage: storage,
		channel: nil,
		conn:    nil,
	}
}

func (p *Scheduler) Connect() error {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s", p.conf.User, p.conf.Pass, p.conf.Host, p.conf.Port)
	slog.Info("connecting to rabbitmq", slog.String("url", fmt.Sprintf("%s:%s", p.conf.Host, p.conf.Port)))

	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	p.conn = conn

	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	p.channel = channel
	err = channel.ExchangeDeclare(
		p.conf.ExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = channel.QueueDeclare(
		p.conf.QueueName, // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	if err != nil {
		return err
	}

	err = channel.QueueBind(
		p.conf.QueueName,
		p.conf.RoutingKey,
		p.conf.ExchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *Scheduler) Close() error {
	slog.Info("Closing connection to AMQP")
	err := p.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *Scheduler) produce(ctx context.Context, event []entities.Event) error {
	for _, event := range event {
		err := p.produceOne(ctx, event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Scheduler) produceOne(ctx context.Context, event entities.Event) error {
	slog.Info("Sending event to queue")
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = p.channel.PublishWithContext(ctx,
		p.conf.ExchangeName,
		p.conf.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return err
}

func (p *Scheduler) DeleteItems() error {
	s := *p.storage
	events, err := s.List()
	if err != nil {
		slog.Error(fmt.Sprintf("Error get list of events: %s", err))
		return err
	}
	for _, event := range events {
		if event.DateTime.Before(time.Now().AddDate(-1, 0, 0)) {
			event := event
			err := s.Delete(&event)
			if err != nil {
				slog.Error(fmt.Sprintf("Error delete event: %s", err))
				return err
			}
		}
	}
	return nil
}

func (p *Scheduler) DatabaseMonitor(ctx context.Context) error {
	var filtered []entities.Event

	s := *p.storage
	events, err := s.List()
	if err != nil {
		return err
	}

	now := time.Now()

	for _, event := range events {
		if !event.Notified {
			notifyTime := event.DateTime.Add(-time.Duration(event.NotifyTime) * time.Second)
			if now.After(notifyTime) || now.Equal(notifyTime) {
				filtered = append(filtered, event)
			}
		}
	}
	err = p.produce(ctx, filtered)
	if err != nil {
		return err
	}
	return nil
}
