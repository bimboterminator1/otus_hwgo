package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

	//nolint:depguard
	storage "github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/storage"
)

type MemoryStorage struct {
	mu     sync.RWMutex
	events map[int64]*storage.Event
	lastID int64
}

func New() storage.Storage {
	return &MemoryStorage{
		events: make(map[int64]*storage.Event),
	}
}

func (m *MemoryStorage) CreateEvent(ctx context.Context, event *storage.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Generate new ID
	m.lastID++
	event.ID = m.lastID

	// Create deep copy to prevent external modifications
	eventCopy := *event
	m.events[event.ID] = &eventCopy

	return nil
}

func (m *MemoryStorage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	existing, exists := m.events[event.ID]
	if !exists {
		return storage.EventNotFound(event.ID)
	}

	// Check if the event belongs to the user
	if existing.UserID != event.UserID {
		return storage.EventNotFound(event.ID)
	}

	// Create deep copy to prevent external modifications
	eventCopy := *event
	m.events[event.ID] = &eventCopy

	return nil
}

func (m *MemoryStorage) DeleteEvent(ctx context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.events[id]; !exists {
		return storage.EventNotFound(id)
	}

	delete(m.events, id)
	return nil
}

func (m *MemoryStorage) GetEvent(ctx context.Context, id int64) (*storage.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	event, exists := m.events[id]
	if !exists {
		return nil, storage.EventNotFound(id)
	}

	// Return copy to prevent external modifications
	eventCopy := *event
	return &eventCopy, nil
}

func (m *MemoryStorage) ListEvents(ctx context.Context, userID int64, from, to time.Time) ([]*storage.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var events []*storage.Event

	for _, event := range m.events {
		if event.UserID == userID &&
			!event.StartTime.After(to) &&
			!event.EndTime.Before(from) {
			// Add copy to prevent external modifications
			eventCopy := *event
			events = append(events, &eventCopy)
		}
	}

	return events, nil
}

func (m *MemoryStorage) ListUpcomingEvents(ctx context.Context, userID int64, limit int) ([]*storage.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	now := time.Now()
	var events []*storage.Event

	// First, collect all upcoming events for the user
	for _, event := range m.events {
		if event.UserID == userID && event.StartTime.After(now) {
			eventCopy := *event
			events = append(events, &eventCopy)
		}
	}

	// Sort events by start time
	sortEventsByStartTime(events)

	// Return only requested number of events
	if len(events) > limit {
		events = events[:limit]
	}

	return events, nil
}

func (m *MemoryStorage) ListEventsNeedingNotification(ctx context.Context, before time.Time) ([]*storage.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var events []*storage.Event

	for _, event := range m.events {
		if !event.NotifyAt.IsZero() && event.NotifyAt.Before(before) {
			eventCopy := *event
			events = append(events, &eventCopy)
		}
	}

	// Sort events by notification time
	sortEventsByNotifyAt(events)

	return events, nil
}

func (m *MemoryStorage) GetEventsByTimeRange(ctx context.Context,
	userID int64,
	start, end time.Time,
) ([]*storage.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var events []*storage.Event

	for _, event := range m.events {
		if event.UserID == userID &&
			(isTimeInRange(event.StartTime, start, end) ||
				isTimeInRange(event.EndTime, start, end) ||
				(event.StartTime.Before(start) && event.EndTime.After(end))) {
			eventCopy := *event
			events = append(events, &eventCopy)
		}
	}

	// Sort events by start time
	sortEventsByStartTime(events)

	return events, nil
}

func (m *MemoryStorage) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clear all data
	m.events = make(map[int64]*storage.Event)
	return nil
}

func (m *MemoryStorage) Ping() error {
	return nil // Memory storage is always available
}

// Helper functions

func sortEventsByStartTime(events []*storage.Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime.Before(events[j].StartTime)
	})
}

func sortEventsByNotifyAt(events []*storage.Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].NotifyAt.Before(events[j].NotifyAt)
	})
}

func isTimeInRange(t, start, end time.Time) bool {
	return !t.Before(start) && !t.After(end)
}
