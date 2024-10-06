package entities

import "time"

type Event struct {
	ID          string
	Title       string
	DateTime    time.Time
	Duration    time.Duration
	Description string
	OwnerId     string
	NotifyTime  int64
}
