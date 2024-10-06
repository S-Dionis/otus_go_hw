package storage

import "github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/entities"

type Storage interface {
	Add(event *entities.Event) error
	Change(event entities.Event) error
	Delete(event entities.Event) error
	List() ([]entities.Event, error)
}
