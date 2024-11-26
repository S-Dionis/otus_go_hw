package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"time"

	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/cmd/config"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/entities"
	_ "github.com/lib/pq" // postgres driver
	"github.com/pressly/goose/v3"
)

type Storage struct {
	connStr string
	db      *sql.DB
}

func New(ctx context.Context, conf config.PsqlConf) *Storage {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Dbname,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("failed to connect to postgres")
		return nil
	}
	err = db.PingContext(ctx)
	if err != nil {
		slog.Error("failed to connect to postgres")
		return nil
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)

	slog.Info(fmt.Sprintf("Run migration with file %s", conf.Migration))
	if err := goose.Up(db, conf.Migration); err != nil {
		log.Fatalf("Ошибка выполнения миграций: %v", err)
	}

	log.Println("Миграции успешно применены!")
	return &Storage{
		connStr: connStr,
		db:      db,
	}
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Add(event *entities.Event) error {
	query := `insert into events(id, title, date_time, duration, description, owner_id, notify_time, notified) 
				values ($1, $2, $3, $4, $5, $6, $7, $8)`

	result, err := s.db.Exec(query, event.ID, event.Title, event.DateTime,
		event.Duration, event.Description, event.OwnerID, event.NotifyTime, event.Notified)
	if err != nil {
		slog.Error(fmt.Sprintf("Error inserting event: %e", err))
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		slog.Error(fmt.Sprintf("Error getting rows affected: %e", err))
		return err
	}

	slog.Info(fmt.Sprintf("affected rows: %s", strconv.FormatInt(affected, 10)))
	return nil
}

func (s *Storage) Change(event *entities.Event) error {
	query := `update events
				  set title = $1, date_time = $2, duration = $3, description = $4, owner_id = $5, notify_time= $6, notified = $7
				  where ID = $8`

	result, err := s.db.Exec(query, event.Title, event.DateTime, event.Duration, event.Description,
		event.OwnerID, event.NotifyTime, event.Notified, event.ID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("affected rows: %s", strconv.FormatInt(affected, 10)))
	return nil
}

func (s *Storage) Delete(event *entities.Event) error {
	query := `DELETE FROM events WHERE events.ID = $1;`

	result, err := s.db.Exec(query, event.ID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("affected rows: %s", strconv.FormatInt(affected, 10)))
	return nil
}

func (s *Storage) List() ([]entities.Event, error) {
	query := `Select id, title, date_time, duration, description, owner_id, notify_time, notified FROM events`

	result, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	events := make([]entities.Event, 0)

	defer result.Close()

	for result.Next() {
		var id string
		var title string
		var dateTime time.Time
		var duration time.Duration
		var description string
		var ownerID string
		var notifyTime int64
		var notified bool

		if err := result.Scan(&id, &title, &dateTime, &duration, &description, &ownerID, &notifyTime, &notified); err != nil {
			return nil, err
		}

		event := entities.Event{
			ID: id, Title: title, DateTime: dateTime, Duration: duration,
			Description: description, OwnerID: ownerID, NotifyTime: notifyTime,
			Notified: notified,
		}
		events = append(events, event)
	}

	return events, nil
}
