package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	app "github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/app"
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/logger"
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/storage"
	"github.com/gorilla/mux"
)

type EventHandlers struct {
	app    *app.App
	logger *logger.Logger
}

func NewEventHandlers(app *app.App, logger *logger.Logger) *EventHandlers {
	return &EventHandlers{
		app:    app,
		logger: logger,
	}
}

func (h *EventHandlers) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
	var event storage.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.logger.Error(fmt.Sprintf("Failed to decode request: %v", err))
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	resEvent, err := h.app.CreateEvent(r.Context(), &event)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to create event: %v", err))
		http.Error(w, "Failed to create event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resEvent)
}

func (h *EventHandlers) HandleUpdateEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var event storage.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	resEvent, err := h.app.UpdateEvent(r.Context(), id, &event)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to update event: %v", err))
		http.Error(w, "Failed to update event", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resEvent)
}

func (h *EventHandlers) HandleDeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := h.app.DeleteEvent(r.Context(), id); err != nil {
		h.logger.Error(fmt.Sprintf("Failed to delete event: %v", err))
		http.Error(w, "Failed to delete event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EventHandlers) HandleListEventsForDay(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		http.Error(w, "Invalid date format. Use RFC3339 format (e.g., 2024-03-20T00:00:00Z)", http.StatusBadRequest)
		return
	}

	// Parse user_id
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endTime := startTime.Add(24 * time.Hour)
	events, err := h.app.ListEvents(r.Context(), userID, startTime, endTime)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to list events: %v", err))
		http.Error(w, "Failed to list events", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func (h *EventHandlers) HandleListEventsForWeek(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		http.Error(w, "Invalid date format. Use RFC3339 format (e.g., 2024-03-20T00:00:00Z)", http.StatusBadRequest)
		return
	}

	// Parse user_id
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endTime := startTime.Add(7 * 24 * time.Hour)

	events, err := h.app.ListEvents(r.Context(), userID, startTime, endTime)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to list events: %v", err))
		http.Error(w, "Failed to list events", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}

func (h *EventHandlers) HandleListEventsForMonth(w http.ResponseWriter, r *http.Request) {
	dateStr := r.URL.Query().Get("date")
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		http.Error(w, "Invalid date format. Use RFC3339 format (e.g., 2024-03-20T00:00:00Z)", http.StatusBadRequest)
		return
	}

	// Parse user_id
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endTime := startTime.AddDate(0, 1, 0)

	events, err := h.app.ListEvents(r.Context(), userID, startTime, endTime)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to list events: %v", err))
		http.Error(w, "Failed to list events", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)
}
