package event

import (
	"bookingService/internal/event"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func NewEventRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetActiveEvents(ctx context.Context) ([]event.Event, error) {
	rows, err := r.db.QueryContext(ctx,
		`
		SELECT id, name, starts_at FROM events
		WHERE starts_at > NOW()
		ORDER BY starts_at ASC
		`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []event.Event
	for rows.Next() {
		var model event.Event
		scanErr := rows.Scan(&model.ID, &model.Name, &model.StartsAt)
		if scanErr != nil {
			return nil, scanErr
		}

		events = append(events, model)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return events, nil
}

func (r *Repository) CreateEvent(ctx context.Context, input event.CreateEventInput) (event.Event, error) {
	var model event.Event

	err := r.db.QueryRowContext(
		ctx,
		`
            INSERT INTO events(name, starts_at)
            VALUES ($1, $2)
            RETURNING id, name, starts_at, created_at
        `,
		input.Name,
		input.StartsAt,
	).Scan(
		&model.ID,
		&model.Name,
		&model.StartsAt,
		&model.CreatedAt,
	)

	if err != nil {
		return event.Event{}, fmt.Errorf("create event: %w", err)
	}

	return model, nil
}

func (r *Repository) GetEventByID(ctx context.Context, id int64) (event.Event, error) {
	var model event.Event

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, starts_at, created_at
				FROM events
				WHERE id = $1`,
		id,
	).Scan(
		&model.ID,
		&model.Name,
		&model.StartsAt,
		&model.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return event.Event{}, sql.ErrNoRows
	}

	if err != nil {
		return event.Event{}, fmt.Errorf("failed get event by id: %w", err)
	}

	return model, nil
}
