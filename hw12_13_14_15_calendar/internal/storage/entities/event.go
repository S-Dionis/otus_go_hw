package entities

import "time"

type Event struct {
	ID          string
	Title       string
	DateTime    time.Time
	Duration    time.Duration
	Description string
	OwnerID     string
	NotifyTime  int64
}
