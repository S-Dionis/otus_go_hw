package internalhttp

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var logger *slog.Logger

func init() {
	file, err := os.OpenFile("method.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		slog.Error("Не удалось открыть файл для логирования: %v", err)
	}

	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelError,
	}

	handler := slog.NewTextHandler(file, handlerOpts)
	logger = slog.New(handler)

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	//defer logFile.Close()

	logger = slog.New(slog.NewTextHandler(logFile, nil))
}

func loggingMiddleware(next http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logString := fmt.Sprintf("%s [%s] %s %s %s %d %s",
			r.RemoteAddr,
			time.Now(),
			r.Method,
			r.URL.Path,
			r.Proto,
			time.Since(start),
			r.UserAgent(),
		)
		logger.Info(logString)
	})
}
