package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/cmd/config"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/server/http"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := config.NewConfig(configFile)
	err := logger.InitLogger(config.Logger.Level)
	if err != nil {
		return
	}

	var storage storage.Storage

	switch config.DBType.Type {
	case "memory":
		storage = memorystorage.New()
	case "sql":
		storage = sqlstorage.New()
	}

	calendar := app.New(storage, config)

	server := internalhttp.NewServer(calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			slog.Error("failed to stop http server: " + err.Error())
		}
	}()

	slog.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		slog.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
