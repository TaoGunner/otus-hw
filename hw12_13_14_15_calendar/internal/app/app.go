package app

import (
	"context"
)

type App struct { // TODO
}

type Storage interface { // TODO
}

func New(_ Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(_ context.Context, _, _ string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
