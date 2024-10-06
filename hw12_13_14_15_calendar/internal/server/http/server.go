package internalhttp

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	app    Application
	server *http.Server
}

type Application interface{}

func NewServer(app Application) *Server {
	s := &Server{
		app: app,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.helloHandler)
	server := &http.Server{
		Addr:              ":8080",
		Handler:           loggingMiddleware(mux),
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
	w.Write([]byte("Hello, World!"))
}
