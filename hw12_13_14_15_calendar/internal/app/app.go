package app

import (
	"context"

	//nolint:depguard
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/logger"
	//nolint:depguard
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/storage"
)

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

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) error {
	if err := a.storage.CreateEvent(ctx, event); err != nil {
		a.logger.Error("Failed to create event: " + err.Error())
		return err
	}

	a.logger.Info("Created event: " + event.Title)
	return nil
}
