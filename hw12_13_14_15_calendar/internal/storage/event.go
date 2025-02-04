package storage

import (
	"context"
	"time"
)

type Event struct {
	ID          int64     `json:"id,omitempty" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	UserID      int64     `json:"user_id" db:"user_id"`
	NotifyAt    time.Time `json:"notify_at,omitempty" db:"notify_at"`
}

type Storage interface {
	CreateEvent(ctx context.Context, event *Event) error
	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEvent(ctx context.Context, id int64) (*Event, error)
	ListEvents(ctx context.Context, userID int64, from, to time.Time) ([]*Event, error)
	ListUpcomingEvents(ctx context.Context, userID int64, limit int) ([]*Event, error)
	ListEventsNeedingNotification(ctx context.Context, before time.Time) ([]*Event, error)
	GetEventsByTimeRange(ctx context.Context, userID int64, start, end time.Time) ([]*Event, error)
	Close() error
	Ping() error
}
