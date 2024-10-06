package app

import (
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/entities"
)

type App struct {
	storage storage.Storage
}

func New(storage storage.Storage) *App {
	return &App{
		storage: storage,
	}
}

func (a *App) CreateEvent(id, title string) error {
	return a.storage.Add(&entities.Event{ID: id, Title: title})
}
