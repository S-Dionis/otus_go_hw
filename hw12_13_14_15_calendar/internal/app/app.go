package app

import (
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/cmd/config"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/entities"
)

type App struct {
	config  config.Config
	storage storage.Storage
}

func New(storage storage.Storage, config config.Config) *App {
	return &App{
		config:  config,
		storage: storage,
	}
}

func (a *App) CreateEvent(id, title string) error {
	return a.storage.Add(&entities.Event{ID: id, Title: title})
}

func (a *App) Storage() storage.Storage {
	return a.storage
}

func (a *App) Config() config.Config {
	return a.config
}
