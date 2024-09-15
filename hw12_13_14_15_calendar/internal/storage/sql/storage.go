package sqlstorage

//nolint:revive
import (
	"database/sql"
	"errors"

	"github.com/TaoGunner/otus-hw/hw12_13_14_15_calendar/internal/storage"
	_ "modernc.org/sqlite"
)

var (
	ErrSQLStorageNotConnected     = errors.New("not connected")
	ErrSQLStorageAlreadyConnected = errors.New("already connected")
)

type Storage struct {
	storage.EventStorer
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	s := &Storage{}
	if err := s.Connect(dbPath); err != nil {
		return nil, err
	}

	return &Storage{}, nil
}

func (s *Storage) initTable() error {
	if s.db == nil {
		return ErrSQLStorageNotConnected
	}

	if _, err := s.db.Exec(eventsCommandCreateTable); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Connect(dbPath string) error {
	if s.db != nil {
		return ErrSQLStorageAlreadyConnected
	}

	dbConn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	s.db = dbConn

	return s.initTable()
}

func (s *Storage) Close() error {
	if s.db == nil {
		return ErrSQLStorageNotConnected
	}

	return s.db.Close()
}

// Add implements storage.EventStore.
func (s *Storage) Add(event storage.Event) (*storage.Event, error) {
	if s.db == nil {
		return nil, ErrSQLStorageNotConnected
	}

	sqlRes, err := s.db.Exec(
		eventsCommandAdd,
		event.ID,
		event.Title,
		event.Datetime,
		event.Duration,
		event.Description,
		event.UserID,
		event.AlarmUntil,
	)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := sqlRes.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, storage.ErrEventTimeIsBusy
	}

	return &event, nil
}

// List implements storage.EventStore.
func (s *Storage) List() ([]storage.Event, error) {
	if s.db == nil {
		return nil, ErrSQLStorageNotConnected
	}

	sqlRows, err := s.db.Query(eventsCommandList)
	if err != nil {
		return nil, err
	}
	defer sqlRows.Close()

	var events []storage.Event

	for sqlRows.Next() {
		var e storage.Event
		if err := sqlRows.Scan(
			&e.ID,
			&e.Title,
			&e.Datetime,
			&e.Duration,
			&e.Description,
			&e.UserID,
			&e.AlarmUntil,
		); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	if err = sqlRows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// Remove implements storage.EventStore.
func (s *Storage) Remove(id string) error {
	if s.db == nil {
		return ErrSQLStorageNotConnected
	}

	sqlRes, err := s.db.Exec(eventsCommandRemove, id)
	if err != nil {
		return err
	}

	rowsAffected, err := sqlRes.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return storage.ErrEventNotFound
	}

	return nil
}

// Update implements storage.EventStore.
func (s *Storage) Update(event storage.Event) error {
	if s.db == nil {
		return ErrSQLStorageNotConnected
	}

	sqlRes, err := s.db.Exec(
		eventsCommandUpdate,
		event.Title,
		event.Datetime,
		event.Duration,
		event.Description,
		event.UserID,
		event.AlarmUntil,
		event.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := sqlRes.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return storage.ErrEventNotFound
	}

	return nil
}
