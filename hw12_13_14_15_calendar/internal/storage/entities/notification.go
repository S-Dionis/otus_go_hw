package entities

import "time"

type Notification struct {
	ID    string
	Title string
	Date  time.Time
	User  string
}
