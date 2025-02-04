package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	//nolint:depguard
	"github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/config"
	//nolint:depguard
	storage "github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/internal/storage"
	//nolint:depguard
	migrator "github.com/bimboterminator1/otus_hwgo/hw12_13_14_15_calendar/migrations"
	"github.com/jmoiron/sqlx"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func New(conf config.StorageConfig) (storage.Storage, error) {
	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=disable",
		conf.Username,
		conf.Password,
		conf.Address,
		conf.Database,
	)

	mg, err := migrator.New(connStr)
	if err != nil {
		return nil, storage.DatabaseError("migration", err)
	}
	defer mg.Close()

	// Run migrations
	if err := mg.Up(); err != nil {
		return nil, storage.DatabaseError("migration", err)
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, storage.DatabaseError("connect", err)
	}

	if conf.Pool.MaxOpenConns > 0 {
		db.SetMaxOpenConns(conf.Pool.MaxOpenConns)
	} else {
		db.SetMaxOpenConns(25) // default
	}

	if conf.Pool.MaxIdleConns > 0 {
		db.SetMaxIdleConns(conf.Pool.MaxIdleConns)
	} else {
		db.SetMaxIdleConns(5) // default
	}
	if conf.Pool.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(conf.Pool.ConnMaxLifetime) * time.Second)
	} else {
		db.SetConnMaxLifetime(5 * time.Minute) // default
	}

	return &PostgresStorage{db: db}, nil
}

func (p *PostgresStorage) CreateEvent(ctx context.Context, event *storage.Event) error {
	const query = `
    INSERT INTO events (
        title, description, start_time, end_time, user_id, notify_at
    ) VALUES (
        :title, :description, :start_time, :end_time, :user_id, :notify_at
    ) RETURNING id`

	rows, err := p.db.NamedQueryContext(ctx, query, event)
	if err != nil {
		return storage.DatabaseError("create", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&event.ID); err != nil {
			return storage.DatabaseError("scan id", err)
		}
	}

	return nil
}

func (p *PostgresStorage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	const query = `
    UPDATE events SET 
        title = :title,
        description = :description,
        start_time = :start_time,
        end_time = :end_time,
        notify_at = :notify_at
    WHERE id = :id AND user_id = :user_id
    `

	result, err := p.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return storage.DatabaseError("update", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return storage.DatabaseError("get affected rows", err)
	}

	if rows == 0 {
		return storage.EventNotFound(event.ID)
	}

	return nil
}

func (p *PostgresStorage) DeleteEvent(ctx context.Context, id int64) error {
	const query = `DELETE FROM events WHERE id = $1`

	result, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return storage.DatabaseError("delete", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return storage.DatabaseError("get affected rows", err)
	}

	if rows == 0 {
		return storage.EventNotFound(id)
	}

	return nil
}

func (p *PostgresStorage) GetEvent(ctx context.Context, id int64) (*storage.Event, error) {
	const query = `
    SELECT id, title, description, start_time, end_time, user_id, notify_at
    FROM events
    WHERE id = $1
    `

	event := &storage.Event{}
	err := p.db.GetContext(ctx, event, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.EventNotFound(id)
	}
	if err != nil {
		return nil, storage.DatabaseError("get", err)
	}

	return event, nil
}

func (p *PostgresStorage) ListEvents(ctx context.Context, userID int64, from, to time.Time) ([]*storage.Event, error) {
	const query = `
    SELECT id, title, description, start_time, end_time, user_id, notify_at
    FROM events
    WHERE user_id = $1
    AND start_time <= $3
    AND end_time >= $2
    ORDER BY start_time ASC
    `
	var events []*storage.Event
	err := p.db.SelectContext(ctx, &events, query, userID, from, to)
	if err != nil {
		return nil, storage.DatabaseError("list", err)
	}

	return events, nil
}

func (p *PostgresStorage) ListUpcomingEvents(ctx context.Context, userID int64, limit int) ([]*storage.Event, error) {
	const query = `
    SELECT id, title, description, start_time, end_time, user_id, notify_at
    FROM events
    WHERE user_id = $1
    AND start_time > CURRENT_TIMESTAMP
    ORDER BY start_time ASC
    LIMIT $2
    `

	var events []*storage.Event
	err := p.db.SelectContext(ctx, &events, query, userID, limit)
	if err != nil {
		return nil, storage.DatabaseError("list upcoming", err)
	}

	return events, nil
}

func (p *PostgresStorage) ListEventsNeedingNotification(ctx context.Context,
	before time.Time) ([]*storage.Event, error) {
	const query = `
    SELECT id, title, description, start_time, end_time, user_id, notify_at
    FROM events
    WHERE notify_at <= $1
    AND notify_at IS NOT NULL
    ORDER BY notify_at ASC
    `

	var events []*storage.Event
	err := p.db.SelectContext(ctx, &events, query, before)
	if err != nil {
		return nil, storage.DatabaseError("list notifications", err)
	}

	return events, nil
}

func (p *PostgresStorage) GetEventsByTimeRange(ctx context.Context,
	userID int64,
	start, end time.Time) ([]*storage.Event, error) {
	const query = `
    SELECT id, title, description, start_time, end_time, user_id, notify_at
    FROM events
    WHERE user_id = $1
    AND (
        (start_time BETWEEN $2 AND $3)
        OR (end_time BETWEEN $2 AND $3)
        OR (start_time <= $2 AND end_time >= $3)
    )
    ORDER BY start_time ASC
    `

	var events []*storage.Event
	err := p.db.SelectContext(ctx, &events, query, userID, start, end)
	if err != nil {
		return nil, storage.DatabaseError("get by time range", err)
	}

	return events, nil
}

func (p *PostgresStorage) Close() error {
	if err := p.db.Close(); err != nil {
		return storage.DatabaseError("close", err)
	}
	return nil
}

func (p *PostgresStorage) Ping() error {
	if err := p.db.Ping(); err != nil {
		return storage.DatabaseError("ping", err)
	}
	return nil
}
