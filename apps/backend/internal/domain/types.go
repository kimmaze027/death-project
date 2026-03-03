package domain

import "time"

type EventType string

const (
	EventTypeAlertSent    EventType = "alert_sent"
	EventTypeAlive        EventType = "alive"
	EventTypeSkipSnoozed  EventType = "skip_snoozed"
	EventTypeSkipSleeping EventType = "skip_sleeping"
	EventTypeSkipCharging EventType = "skip_charging"
)

type Event struct {
	EventID    string            `json:"event_id"`
	DeviceID   string            `json:"device_id"`
	EventType  EventType         `json:"event_type"`
	OccurredAt time.Time         `json:"occurred_at"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type Snooze struct {
	DeviceID      string    `json:"device_id"`
	DurationHours int       `json:"duration_hours"`
	StartAt       time.Time `json:"start_at"`
	EndAt         time.Time `json:"end_at"`
}

type LatestStatus struct {
	DeviceID        string    `json:"device_id"`
	LatestEventType EventType `json:"latest_event_type"`
	LatestEventAt   time.Time `json:"latest_event_at"`
	IsSleeping      bool      `json:"is_sleeping"`
	IsCharging      bool      `json:"is_charging"`
	SnoozeUntil     time.Time `json:"snooze_until"`
}
