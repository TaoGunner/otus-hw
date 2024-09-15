package memorystorage

import (
	"errors"
	"sort"
	"sync"

	"github.com/TaoGunner/otus-hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"golang.org/x/exp/maps"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() storage.EventStorer {
	return &Storage{events: map[string]storage.Event{}}
}

// Add добавляет событие.
func (s *Storage) Add(event storage.Event) (*storage.Event, error) {
	// Проверка на то, что время не занято
	if _, err := s.EventAtTime(event.Datetime); !errors.Is(err, storage.ErrEventNotFound) {
		return nil, storage.ErrEventTimeIsBusy
	}

	s.mu.RLock()
	s.events[event.ID] = event
	s.mu.RUnlock()

	return &event, nil
}

// EventAtTime возвращает событие по запрошенной метке времени.
func (s *Storage) EventAtTime(ts int64) (*storage.Event, error) {
	for _, e := range s.events {
		if e.Datetime == ts {
			return &e, nil
		}
	}

	return nil, storage.ErrEventNotFound
}

// List возвращает отсортированный список всех событий.
func (s *Storage) List() ([]storage.Event, error) {
	events := maps.Values(s.events)

	sort.Slice(events, func(i, j int) bool {
		return events[i].Datetime < events[j].Datetime
	})

	return events, nil
}

// Remove удаляет событие с указанным id.
func (s *Storage) Remove(id string) error {
	// Проверка, что событие с ID существует
	if _, ok := s.events[id]; !ok {
		return storage.ErrEventNotFound
	}

	s.mu.RLock()
	delete(s.events, id)
	s.mu.RUnlock()

	return nil
}

// Update обновляет запись о событии.
func (s *Storage) Update(event storage.Event) error {
	// Проверка, что ID события передан
	if _, err := uuid.Parse(event.ID); err != nil {
		return storage.ErrEventInvalidUUID
	}

	// Проверка, что событие с ID существует
	if _, ok := s.events[event.ID]; !ok {
		return storage.ErrEventNotFound
	}

	s.mu.RLock()
	s.events[event.ID] = event
	s.mu.RUnlock()

	return nil
}
