package internalhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newServer() *Server {
	storage := memorystorage.New()
	appConfig := config.Config{
		Logger: config.LoggerConf{Level: "INFO"},
		Server: config.ServerConf{
			Host: "localhost",
			Port: "80",
		},
		DBType: config.DBType{
			Type: "memory",
		},
		GRPCConf: config.GRPCConf{
			Port: "8080",
		},
	}
	application := app.New(storage, appConfig)
	server := NewServer(application)
	return server
}

func TestHelloHandler(t *testing.T) {
	event := entities.Event{
		ID:          "1",
		Title:       "mur",
		DateTime:    time.Now(),
		Duration:    10,
		Description: "mewuk",
		OwnerID:     "1",
		NotifyTime:  0,
	}
	eventJSON, err := json.Marshal(event)
	require.NoError(t, err)

	t.Run("Hello event handler", func(t *testing.T) {
		server := newServer()
		writer := httptest.NewRecorder()
		server.helloHandler(writer, nil)

		response := writer.Result()
		defer response.Body.Close()

		data, err := io.ReadAll(response.Body)
		require.NoError(t, err)

		expected := string(data)
		actual := "Hello, World!"
		require.Equal(t, expected, actual)
	})

	t.Run("Add events handler test", func(t *testing.T) {
		server := newServer()
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(eventJSON))
		recorder := httptest.NewRecorder()
		server.addEventsHandler(recorder, req)
		resp := recorder.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete events handler test", func(t *testing.T) {
		server := newServer()
		req := httptest.NewRequest(http.MethodDelete, "/events", bytes.NewReader(eventJSON))
		recorder := httptest.NewRecorder()
		server.addEventsHandler(recorder, req)
		resp := recorder.Result()
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("List events handler test", func(t *testing.T) {
		server := newServer()
		e := &entities.Event{
			ID:          "1",
			Title:       "mur",
			DateTime:    time.Now(),
			Duration:    10,
			Description: "mewuk",
			OwnerID:     "1",
			NotifyTime:  0,
		}
		err := server.app.Storage().Add(e)
		require.NoError(t, err)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/events?period=DAY", nil)
		server.getEventsHandler(recorder, req)
		resp := recorder.Result()
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		var events []entities.Event
		body := resp.Body
		fmt.Print(body)
		err = json.NewDecoder(body).Decode(&events)
		require.NoError(t, err)
		require.Len(t, events, 1)
		event := events[0]
		require.Equal(t, e.ID, event.ID)
		require.Equal(t, e.Duration, event.Duration)
		require.Equal(t, e.Title, event.Title)
		require.Equal(t, e.NotifyTime, event.NotifyTime)
		require.Equal(t, e.Description, event.Description)
		require.Equal(t, e.OwnerID, event.OwnerID)
	})
}
