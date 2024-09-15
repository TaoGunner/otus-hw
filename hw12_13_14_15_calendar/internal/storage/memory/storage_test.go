package memorystorage

import (
	"testing"
	"time"

	"github.com/TaoGunner/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

func TestStorage(t *testing.T) {
	var err error
	db := New()
	event := &storage.Event{
		ID:          uuid.NewString(),
		Title:       "Тестовое событие",
		Datetime:    time.Now().Unix(),
		Duration:    3600,
		Description: "",
		UserID:      "",
		AlarmUntil:  0,
	}

	t.Run("Storage CRUD", func(t *testing.T) {
		// Добавление события
		event, err = db.Add(*event)
		if err != nil {
			t.Error(err)
		}

		// Добавление события на занятое время
		if _, err := db.Add(*event); err == nil {
			t.Error("must be 'time is busy' error")
		}

		// Обновление существующего события
		event.Description = "test description"
		if err := db.Update(*event); err != nil {
			t.Error(err)
		}

		// Попытка обновления несуществующего события
		event.ID = uuid.New().String()
		if err := db.Update(*event); err == nil {
			t.Error("must be 'event not found' error")
		}

		// Получение списка событий (размер 1)
		eventList, err := db.List()
		if err != nil {
			t.Error(err)
		}
		if len(eventList) != 1 {
			t.Error("event list length must be 1")
		}

		// Удаление существующего события
		event = &eventList[0]
		if err := db.Remove(event.ID); err != nil {
			t.Error(err)
		}

		// Попытка удаление несуществующего события
		if err := db.Remove(event.ID); err == nil {
			t.Error("must be 'event not found' error")
		}
	})
}
