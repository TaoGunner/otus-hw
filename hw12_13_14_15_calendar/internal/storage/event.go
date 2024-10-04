package storage

import "errors"

var (
	ErrEventInvalidUUID = errors.New("empty uuid")
	ErrEventNotFound    = errors.New("event not found")
	ErrEventTimeIsBusy  = errors.New("time is busy")
)

type EventStorer interface {
	Add(event Event) (*Event, error)
	Update(event Event) error
	Remove(id string) error
	List() ([]Event, error)
	EventAtTime(ts int64) (*Event, error)
}

type Event struct {
	ID          string
	Title       string
	Datetime    int64
	Duration    int64
	Description string
	UserID      string
	AlarmUntil  int64
}
