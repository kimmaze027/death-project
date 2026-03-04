package store

import (
	"errors"
	"sync"
	"time"

	"github.com/kimmaze027/death-project/apps/backend/internal/domain"
)

type MemoryStore struct {
	mu       sync.RWMutex
	events   map[string][]domain.Event
	snoozes  map[string]domain.Snooze
	statuses map[string]domain.LatestStatus
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		events:   make(map[string][]domain.Event),
		snoozes:  make(map[string]domain.Snooze),
		statuses: make(map[string]domain.LatestStatus),
	}
}

func (m *MemoryStore) AddEvent(event domain.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.events[event.DeviceID] = append(m.events[event.DeviceID], event)

	status := domain.LatestStatus{
		DeviceID:        event.DeviceID,
		LatestEventType: event.EventType,
		LatestEventAt:   event.OccurredAt,
		IsSleeping:      event.EventType == domain.EventTypeSkipSleeping,
		IsCharging:      event.EventType == domain.EventTypeSkipCharging,
	}

	if snooze, ok := m.snoozes[event.DeviceID]; ok && snooze.EndAt.After(time.Now()) {
		status.SnoozeUntil = snooze.EndAt
	}

	m.statuses[event.DeviceID] = status
}

func (m *MemoryStore) SetSnooze(s domain.Snooze) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.snoozes[s.DeviceID] = s
}

func (m *MemoryStore) LatestStatus(deviceID string) (domain.LatestStatus, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	status, ok := m.statuses[deviceID]
	if !ok {
		return domain.LatestStatus{}, errors.New("status not found")
	}
	return status, nil
}
