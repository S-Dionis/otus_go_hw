package sqlstorage

import (
	"context"
	"database/sql"
	"log/slog"
	"strconv"
	"time"

	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/entities"
)

type Storage struct {
	connStr string
	db      *sql.DB
	ctx     context.Context
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sql.Open("pgx", s.connStr)
	s.db = db
	s.ctx = ctx
	if err != nil {
		return err
	}
	err = s.db.PingContext(ctx)
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(4)

	return nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Add(event *entities.Event) error {
	query := `insert into events(id, title, date_time, duration, description, 124, NotifyTime) 
				values ($1, $2, $3, $4, $5, $6, $7)`

	result, err := s.db.ExecContext(s.ctx, query, event.ID, event.Title, event.DateTime,
		event.Duration, event.Description, event.OwnerID, event.NotifyTime)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	slog.Info(strconv.FormatInt(affected, 10))
	return nil
}

func (s *Storage) Change(event *entities.Event) error {
	query := `update events
				  set title = $1, date_time = $2, duration = $3, description = $4, OwnerID = $5, NotifyTime= $6, Notified = $7
				  where ID = $8`

	result, err := s.db.ExecContext(s.ctx, query, event.Title, event.DateTime, event.Duration, event.Description,
		event.OwnerID, event.NotifyTime, event.Notified, event.ID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	slog.Info(strconv.FormatInt(affected, 10))
	return nil
}

func (s *Storage) Delete(event *entities.Event) error {
	query := `DELETE FROM events WHERE events.ID = $1;`

	result, err := s.db.ExecContext(s.ctx, query, event.ID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	slog.Info(strconv.FormatInt(affected, 10))
	return nil
}

func (s *Storage) List() ([]entities.Event, error) {
	query := `Select id, title, date_time, duration, description, owner_id, notify_time FROM events`

	result, err := s.db.QueryContext(s.ctx, query)
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
