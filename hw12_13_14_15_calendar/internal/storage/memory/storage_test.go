package memorystorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	storage "github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestMemoryStorage_CreateEvent(t *testing.T) {
	store := New()
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		event := &storage.Event{
			Title:     "Test Event",
			UserID:    1,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(time.Hour),
		}

		err := store.CreateEvent(ctx, event)
		require.NoError(t, err)
		require.Greater(t, event.ID, int64(0))

		// Verify event was stored
		stored, err := store.GetEvent(ctx, event.ID)
		require.NoError(t, err)
		require.Equal(t, event.Title, stored.Title)
	})
}

func TestMemoryStorage_UpdateEvent(t *testing.T) {
	store := New()
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		// Create initial event
		event := &storage.Event{
			Title:     "Initial Title",
			UserID:    1,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(time.Hour),
		}
		require.NoError(t, store.CreateEvent(ctx, event))

		// Update event
		event.Title = "Updated Title"
		err := store.UpdateEvent(ctx, event)
		require.NoError(t, err)

		// Verify update
		updated, err := store.GetEvent(ctx, event.ID)
		require.NoError(t, err)
		require.Equal(t, "Updated Title", updated.Title)
	})

	t.Run("update non-existent event", func(t *testing.T) {
		event := &storage.Event{
			ID:     999,
			UserID: 1,
			Title:  "Non-existent",
		}
		err := store.UpdateEvent(ctx, event)
		require.Error(t, err)
	})

	t.Run("update event with wrong user", func(t *testing.T) {
		// Create event
		event := &storage.Event{
			Title:     "Original",
			UserID:    1,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(time.Hour),
		}
		require.NoError(t, store.CreateEvent(ctx, event))

		// Try to update with different user
		event.UserID = 2
		err := store.UpdateEvent(ctx, event)
		require.Error(t, err)
	})
}

func TestMemoryStorage_DeleteEvent(t *testing.T) {
	store := New()
	ctx := context.Background()

	t.Run("successful deletion", func(t *testing.T) {
		event := &storage.Event{
			Title:     "To Delete",
			UserID:    1,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(time.Hour),
		}
		require.NoError(t, store.CreateEvent(ctx, event))

		err := store.DeleteEvent(ctx, event.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = store.GetEvent(ctx, event.ID)
		require.Error(t, err)
	})

	t.Run("delete non-existent event", func(t *testing.T) {
		err := store.DeleteEvent(ctx, 999)
		require.Error(t, err)
	})
}

func TestMemoryStorage_ListEvents(t *testing.T) {
	store := New()
	ctx := context.Background()
	now := time.Now()

	// Create test events
	events := []*storage.Event{
		{
			Title:     "Event 1",
			UserID:    1,
			StartTime: now,
			EndTime:   now.Add(time.Hour),
		},
		{
			Title:     "Event 2",
			UserID:    1,
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(3 * time.Hour),
		},
		{
			Title:     "Other User Event",
			UserID:    2,
			StartTime: now,
			EndTime:   now.Add(time.Hour),
		},
	}

	for _, e := range events {
		require.NoError(t, store.CreateEvent(ctx, e))
	}

	t.Run("list user events in time range", func(t *testing.T) {
		listed, err := store.ListEvents(ctx, 1, now, now.Add(4*time.Hour))
		require.NoError(t, err)
		require.Len(t, listed, 2)
	})

	t.Run("list events with empty result", func(t *testing.T) {
		listed, err := store.ListEvents(ctx, 1, now.Add(5*time.Hour), now.Add(6*time.Hour))
		require.NoError(t, err)
		require.Empty(t, listed)
	})
}

func TestMemoryStorage_ListUpcomingEvents(t *testing.T) {
	store := New()
	ctx := context.Background()
	now := time.Now()

	// Create test events
	events := []*storage.Event{
		{
			Title:     "Past Event",
			UserID:    1,
			StartTime: now.Add(-2 * time.Hour),
			EndTime:   now.Add(-1 * time.Hour),
		},
		{
			Title:     "Future Event 1",
			UserID:    1,
			StartTime: now.Add(time.Hour),
			EndTime:   now.Add(2 * time.Hour),
		},
		{
			Title:     "Future Event 2",
			UserID:    1,
			StartTime: now.Add(3 * time.Hour),
			EndTime:   now.Add(4 * time.Hour),
		},
	}

	for _, e := range events {
		require.NoError(t, store.CreateEvent(ctx, e))
	}

	t.Run("list upcoming events with limit", func(t *testing.T) {
		listed, err := store.ListUpcomingEvents(ctx, 1, 1)
		require.NoError(t, err)
		require.Len(t, listed, 1)
		require.Equal(t, "Future Event 1", listed[0].Title)
	})
}

func TestMemoryStorage_ListEventsNeedingNotification(t *testing.T) {
	store := New()
	ctx := context.Background()
	now := time.Now()

	events := []*storage.Event{
		{
			Title:     "Notify Soon",
			UserID:    1,
			StartTime: now.Add(time.Hour),
			EndTime:   now.Add(2 * time.Hour),
			NotifyAt:  now.Add(15 * time.Minute),
		},
		{
			Title:     "Notify Later",
			UserID:    1,
			StartTime: now.Add(3 * time.Hour),
			EndTime:   now.Add(4 * time.Hour),
			NotifyAt:  now.Add(2 * time.Hour),
		},
	}

	for _, e := range events {
		require.NoError(t, store.CreateEvent(ctx, e))
	}

	t.Run("list events needing notification", func(t *testing.T) {
		listed, err := store.ListEventsNeedingNotification(ctx, now.Add(time.Hour))
		require.NoError(t, err)
		require.Len(t, listed, 1)
		require.Equal(t, "Notify Soon", listed[0].Title)
	})
}

func TestMemoryStorage_GetEventsByTimeRange(t *testing.T) {
	store := New()
	ctx := context.Background()
	now := time.Now()

	events := []*storage.Event{
		{
			Title:     "Event In Range",
			UserID:    1,
			StartTime: now.Add(time.Hour),
			EndTime:   now.Add(2 * time.Hour),
		},
		{
			Title:     "Event Outside Range",
			UserID:    1,
			StartTime: now.Add(4 * time.Hour),
			EndTime:   now.Add(5 * time.Hour),
		},
	}

	for _, e := range events {
		require.NoError(t, store.CreateEvent(ctx, e))
	}

	t.Run("get events in time range", func(t *testing.T) {
		listed, err := store.GetEventsByTimeRange(ctx, 1, now, now.Add(3*time.Hour))
		require.NoError(t, err)
		require.Len(t, listed, 1)
		require.Equal(t, "Event In Range", listed[0].Title)
	})
}

func TestMemoryStorage_Concurrent(t *testing.T) {
	store := New()
	ctx := context.Background()
	const numGoroutines = 10

	t.Run("concurrent creation", func(t *testing.T) {
		done := make(chan bool)
		for i := 0; i < numGoroutines; i++ {
			go func(i int) {
				event := &storage.Event{
					Title:     fmt.Sprintf("Concurrent Event %d", i),
					UserID:    1,
					StartTime: time.Now(),
					EndTime:   time.Now().Add(time.Hour),
				}
				err := store.CreateEvent(ctx, event)
				require.NoError(t, err)
				done <- true
			}(i)
		}

		for i := 0; i < numGoroutines; i++ {
			<-done
		}
	})
}

func TestMemoryStorage_CloseAndPing(t *testing.T) {
	store := New()

	t.Run("close storage", func(t *testing.T) {
		err := store.Close()
		require.NoError(t, err)
	})

	t.Run("ping storage", func(t *testing.T) {
		err := store.Ping()
		require.NoError(t, err)
	})
}
