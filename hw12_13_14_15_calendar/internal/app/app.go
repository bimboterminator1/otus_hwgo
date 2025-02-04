package app

import (
	"context"
	"fmt"
	"time"

	//nolint:depguard
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/logger"
	//nolint:depguard
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/storage"
)

type Server interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type App struct {
	logger  logger.Logger
	storage storage.Storage
}

func New(logger logger.Logger, storage storage.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) (*storage.Event, error) {
	if err := validateEvent(event); err != nil {
		a.logger.Error(fmt.Sprintf("Event validation failed %v", err))
		return nil, fmt.Errorf("event validation failed: %w", err)
	}

	if err := a.storage.CreateEvent(ctx, event); err != nil {
		a.logger.Error(fmt.Sprintf("Failed to create event: %v", err))
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// Get the created event to return it with the assigned ID
	createdEvent, err := a.storage.GetEvent(ctx, event.ID)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to get created event: %v", err))
		return nil, fmt.Errorf("failed to get created event: %w", err)
	}

	a.logger.Info(fmt.Sprintf("Created event successfully %s", event.Title))
	return createdEvent, nil
}

func (a *App) UpdateEvent(ctx context.Context, id int64, event *storage.Event) (*storage.Event, error) {
	if err := validateEvent(event); err != nil {
		a.logger.Error(fmt.Sprintf("Event validation failed %v", err))
		return nil, fmt.Errorf("event validation failed: %w", err)
	}

	// Check if event exists
	existing, err := a.storage.GetEvent(ctx, id)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to get existing event: %v", err))
		return nil, fmt.Errorf("failed to get existing event: %w", err)
	}

	if existing.UserID != event.UserID {
		a.logger.Error(fmt.Sprintf("Unauthorized event update attempt: %v", err))
		return nil, fmt.Errorf("unauthorized: event belongs to different user")
	}

	event.ID = id // Ensure the ID is set correctly
	if err := a.storage.UpdateEvent(ctx, event); err != nil {
		a.logger.Error(fmt.Sprintf("Failed to update event: %v", err))
		return nil, fmt.Errorf("failed to update event: %w", err)
	}

	// Get the updated event to return
	updatedEvent, err := a.storage.GetEvent(ctx, id)
	if err != nil {
		a.logger.Error(fmt.Sprintf("Failed to get updated event: %v", err))
		return nil, fmt.Errorf("failed to get updated event: %w", err)
	}

	a.logger.Info(fmt.Sprintf("Updated event successfully %d", id))
	return updatedEvent, nil
}

func (a *App) DeleteEvent(ctx context.Context, id int64) error {
	if err := a.storage.DeleteEvent(ctx, id); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	a.logger.Info(fmt.Sprintf("Deleted event successfully %d", id))
	return nil
}

func (a *App) ListEvents(ctx context.Context, userID int64, startTime, endTime time.Time) ([]*storage.Event, error) {
	if !startTime.Before(endTime) {
		return nil, fmt.Errorf("invalid time range: start time must be before end time")
	}

	events, err := a.storage.ListEvents(ctx, userID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	return events, nil
}

// Helper functions.

func validateEvent(event *storage.Event) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	if event.Title == "" {
		return fmt.Errorf("event title cannot be empty")
	}

	if event.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	if event.StartTime.IsZero() {
		return fmt.Errorf("start time cannot be zero")
	}

	if event.EndTime.IsZero() {
		return fmt.Errorf("end time cannot be zero")
	}

	if !event.StartTime.Before(event.EndTime) {
		return fmt.Errorf("end time must be after start time")
	}

	if !event.NotifyAt.IsZero() && !event.NotifyAt.Before(event.StartTime) {
		return fmt.Errorf("notification time must be before start time")
	}

	return nil
}
