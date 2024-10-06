package internalhttp

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	app Application
}

type Application interface {
}

func NewServer(app Application) *Server {
	return &Server{
		app: app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.helloHandler)
	server := &http.Server{
		Addr:    ":8080",
		Handler: loggingMiddleware(mux),
	}

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			slog.Info("server" + string(time.Now().Format(time.RFC3339)) + "Shutdown")
			slog.Error("Server Shutdown Error:", err)
		}
	}()

	slog.Info("server" + time.Now().Format(time.RFC3339) + "Start")
	return server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.Stop(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
