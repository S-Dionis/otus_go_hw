package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/entities"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

const (
	serverURL = "http://calendar:8888"
	dbConnStr = "host=postgres port=5432 user=user password=password dbname=calendar sslmode=disable"
)

func resetSubSeconds(t time.Time) time.Time {
	return time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(),
		0, t.Location(),
	)
}

func TestIntegration(t *testing.T) {
	<-time.After(5 * time.Second)

	fmt.Println("Run integration tests...")

	db, err := sql.Open("postgres", dbConnStr)
	require.NoError(t, err)
	deleteQuery := "DELETE FROM events"
	_, err = db.Exec(deleteQuery)
	require.NoError(t, err)

	event := entities.Event{
		ID:          "1",
		Title:       "Test Event",
		DateTime:    time.Now().Add(time.Second * 10),
		Duration:    time.Hour,
		Description: "This is a test event",
		OwnerID:     "owner123",
		NotifyTime:  time.Now().Add(-time.Second * 3).Unix(),
		Notified:    false,
	}

	postBody, _ := json.Marshal(event)
	addEventURL := fmt.Sprintf("%s/events", serverURL)
	fmt.Println("run request to url:", addEventURL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, addEventURL, bytes.NewBuffer(postBody))
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var dbEvent entities.Event
	query := `SELECT id, title, date_time, duration, description, owner_id, notify_time, notified 
			  FROM events`
	row := db.QueryRow(query)
	err = row.Scan(&dbEvent.ID, &dbEvent.Title, &dbEvent.DateTime, &dbEvent.Duration, &dbEvent.Description,
		&dbEvent.OwnerID, &dbEvent.NotifyTime, &dbEvent.Notified)
	require.NoError(t, err)
	dbEvent.DateTime = dbEvent.DateTime.In(time.Local)
	dbEvent.DateTime = resetSubSeconds(dbEvent.DateTime)
	event.DateTime = resetSubSeconds(event.DateTime)
	require.Equal(t, event, dbEvent)

	slog.Info("wait for notification")
	<-time.After(20 * time.Second)
	slog.Info("stop waiting and check")
	row = db.QueryRow(query)
	err = row.Scan(&dbEvent.ID, &dbEvent.Title, &dbEvent.DateTime, &dbEvent.Duration, &dbEvent.Description,
		&dbEvent.OwnerID, &dbEvent.NotifyTime, &dbEvent.Notified)
	require.NoError(t, err)
	require.Equal(t, true, dbEvent.Notified)
}
