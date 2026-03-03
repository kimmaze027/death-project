package store

import (
	"testing"
	"time"

	"github.com/kimmaze027/death-project/apps/backend/internal/domain"
)

func TestMemoryStore_LatestStatus_NotFound(t *testing.T) {
	t.Parallel()

	s := NewMemoryStore()
	_, err := s.LatestStatus("missing-device")
	if err == nil {
		t.Fatal("expected error for missing status")
	}
}

func TestMemoryStore_AddEvent_UpdatesLatestStatus(t *testing.T) {
	t.Parallel()

	s := NewMemoryStore()
	at := time.Date(2026, 3, 3, 10, 0, 0, 0, time.UTC)

	s.AddEvent(domain.Event{
		EventID:    "evt-1",
		DeviceID:   "watch-1",
		EventType:  domain.EventTypeAlive,
		OccurredAt: at,
	})

	status, err := s.LatestStatus("watch-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.DeviceID != "watch-1" {
		t.Fatalf("unexpected device id: %s", status.DeviceID)
	}
	if status.LatestEventType != domain.EventTypeAlive {
		t.Fatalf("unexpected event type: %s", status.LatestEventType)
	}
	if !status.LatestEventAt.Equal(at) {
		t.Fatalf("unexpected occurred_at: %s", status.LatestEventAt)
	}
	if status.IsSleeping || status.IsCharging {
		t.Fatal("alive event should not set sleeping/charging flags")
	}
}

func TestMemoryStore_AddEvent_SetsSkipFlags(t *testing.T) {
	t.Parallel()

	s := NewMemoryStore()
	at := time.Date(2026, 3, 3, 11, 0, 0, 0, time.UTC)

	s.AddEvent(domain.Event{
		EventID:    "evt-2",
		DeviceID:   "watch-2",
		EventType:  domain.EventTypeSkipSleeping,
		OccurredAt: at,
	})

	status, err := s.LatestStatus("watch-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.IsSleeping {
		t.Fatal("skip_sleeping should set sleeping=true")
	}
	if status.IsCharging {
		t.Fatal("skip_sleeping should set charging=false")
	}
}

func TestMemoryStore_SetSnooze_ReflectedInLatestStatus(t *testing.T) {
	t.Parallel()

	s := NewMemoryStore()
	futureEnd := time.Now().Add(2 * time.Hour)

	s.SetSnooze(domain.Snooze{
		DeviceID:      "watch-3",
		DurationHours: 2,
		StartAt:       time.Now(),
		EndAt:         futureEnd,
	})

	s.AddEvent(domain.Event{
		EventID:    "evt-3",
		DeviceID:   "watch-3",
		EventType:  domain.EventTypeAlertSent,
		OccurredAt: time.Now(),
	})

	status, err := s.LatestStatus("watch-3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.SnoozeUntil.IsZero() {
		t.Fatal("expected snooze_until to be set")
	}
	if !status.SnoozeUntil.Equal(futureEnd) {
		t.Fatalf("unexpected snooze_until: %s", status.SnoozeUntil)
	}
}
