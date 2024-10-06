package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	level string
}

func InitLogger(level string) error {

	var slogLevel slog.Level
	var err = slogLevel.UnmarshalText([]byte(level))

	if err == nil {
		return err
	}

	handlerOpts := &slog.HandlerOptions{
		Level: slogLevel,
	}

	logHandler := slog.NewTextHandler(os.Stdout, handlerOpts)

	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	slog.Info("Logger slog successfully initialized")
	return nil
}
