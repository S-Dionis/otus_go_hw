package memorystorage

import (
	"errors"
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/entities"
	"github.com/google/uuid"
	"sync"
)

type Storage struct {
	mu           sync.RWMutex //nolint:unused
	events       map[string]*entities.Event
	lastInsertId string
}

func New() *Storage {
	return &Storage{
		events: make(map[string]*entities.Event),
	}
}

func (s *Storage) Add(event *entities.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	event.ID = id
	s.events[id] = event

	return nil
}

func (s *Storage) Change(event entities.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	e := s.events[event.ID]

	if e == nil {
		return errors.New("event not found")
	}

	e.Title = event.Title
	e.DateTime = event.DateTime
	e.Duration = event.Duration
	e.Description = event.Description
	e.OwnerId = event.OwnerId
	e.NotifyTime = event.NotifyTime

	return nil
}

func (s *Storage) Delete(event entities.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, event.ID)
	return nil
}

func (s *Storage) List() ([]entities.Event, error) {
	arr := make([]entities.Event, 0)

	for _, v := range s.events {
		arr = append(arr, *v)
	}

	return arr, nil
}
