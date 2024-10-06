package memorystorage

import (
	"github.com/S-Dionis/otus_go_hw/hw12_13_14_15_calendar/internal/storage/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	storage := New()
	t.Run("Test memory storage", func(t *testing.T) {
		event := entities.Event{
			ID:          "",
			Title:       "Test title",
			DateTime:    time.Time{},
			Duration:    0,
			Description: "Test description",
			OwnerId:     "owner id",
			NotifyTime:  0,
		}

		require.NoError(t, storage.Add(&event))

		list, err := storage.List()

		require.NoError(t, err)
		require.Len(t, list, 1)
		assert.Equal(t, event, list[0])

		eventChanged := entities.Event{
			ID:          event.ID,
			Title:       "Test title changed",
			DateTime:    time.Time{},
			Duration:    0,
			Description: "Test description changed",
			OwnerId:     "owner id",
			NotifyTime:  0,
		}

		require.NoError(t, storage.Change(eventChanged))

		listChanged, err := storage.List()
		require.NoError(t, err)
		require.Len(t, listChanged, 1)
		assert.Equal(t, eventChanged, listChanged[0])

	})
}
