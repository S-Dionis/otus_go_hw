package internalhttp

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/entities"
	"github.com/gorilla/mux"
)

type Server struct {
	app    *app.App
	server *http.Server
}

func NewServer(app *app.App) *Server {
	serverMux := http.NewServeMux()
	s := &Server{
		app: app,
	}
	serverMux.HandleFunc("/hello", s.helloHandler)
	serverMux.HandleFunc("/events", s.eventsHandler)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           loggingMiddleware(serverMux),
		ReadHeaderTimeout: 2 * time.Second,
	}
	s.server = server

	return s
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		if err := s.server.Shutdown(context.Background()); err != nil {
			slog.Info("server" + time.Now().Format(time.RFC3339) + "Shutdown")
			slog.Error("Server Shutdown Error:", err)
		}
	}()

	slog.Info("server" + time.Now().Format(time.RFC3339) + "Start")
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) helloHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Hello, World!"))
	if err != nil {
		return
	}
}

func (s *Server) eventsHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		s.getEventsHandler(w, req)
	case http.MethodPost:
		s.addEventsHandler(w, req)
	case http.MethodPut:
		s.updateEventsHandler(w, req)
	case http.MethodDelete:
		s.deleteEventsHandler(w, req)
	}
}

func getEvent(req *http.Request) (*entities.Event, error) {
	var event entities.Event
	if err := json.NewDecoder(req.Body).Decode(&event); err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *Server) deleteEventsHandler(w http.ResponseWriter, req *http.Request) {
	event, err := getEvent(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.Storage().Delete(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) updateEventsHandler(w http.ResponseWriter, req *http.Request) {
	event, err := getEvent(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.Storage().Change(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) addEventsHandler(w http.ResponseWriter, req *http.Request) {
	event, err := getEvent(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.Storage().Add(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) getEventsHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	period := vars["period"]

	events, err := s.app.Storage().List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	switch period {
	case "DAY":
		slog.Info("Get day events")
		events = entities.GetTodayEvents(events)
	case "WEEK":
		slog.Info("Get weekly events")
		events = entities.GetWeekEvents(events)
	case "MONTH":
		slog.Info("Get monthly events")
		events = entities.GetMonthEvents(events)
	default:
		slog.Info("Get all events")
	}
	if err := json.NewEncoder(w).Encode(events); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
