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
	Notified    bool
}

func GetTodayEvents(events []Event) []Event {
	var todayEvents []Event

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	for _, event := range events {
		if event.DateTime.After(startOfDay) && event.DateTime.Before(endOfDay) {
			todayEvents = append(todayEvents, event)
		}
	}
	return todayEvents
}

func GetWeekEvents(events []Event) []Event {
	var weekEvents []Event
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()-time.Monday))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	for _, event := range events {
		if event.DateTime.After(startOfWeek) && event.DateTime.Before(endOfWeek) {
			weekEvents = append(weekEvents, event)
		}
	}
	return weekEvents
}

func GetMonthEvents(events []Event) []Event {
	var monthEvents []Event
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	for _, event := range events {
		if event.DateTime.After(startOfMonth) && event.DateTime.Before(endOfMonth) {
			monthEvents = append(monthEvents, event)
		}
	}
	return monthEvents
}
