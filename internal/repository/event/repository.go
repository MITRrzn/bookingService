package event

import (
	"bookingService/internal/structs"
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

func (r *Repository) GetActiveEvents(ctx context.Context) ([]structs.Event, error) {
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

	var events []structs.Event
	for rows.Next() {
		var event structs.Event
		scanErr := rows.Scan(&event.ID, &event.Name, &event.StartsAt)
		if scanErr != nil {
			return nil, scanErr
		}

		events = append(events, event)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return events, nil
}

func (r *Repository) CreateEvent(ctx context.Context, input structs.CreateEventInput) (structs.Event, error) {
	var event structs.Event

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
		&event.ID,
		&event.Name,
		&event.StartsAt,
		&event.CreatedAt,
	)

	if err != nil {
		return structs.Event{}, fmt.Errorf("create event: %w", err)
	}

	return event, nil
}

func (r *Repository) GetEventByID(ctx context.Context, id int64) (structs.Event, error) {
	var event structs.Event

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, starts_at, created_at
				FROM events
				WHERE id = $1`,
		id,
	).Scan(
		&event.ID,
		&event.Name,
		&event.StartsAt,
		&event.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return structs.Event{}, sql.ErrNoRows
	}

	if err != nil {
		return structs.Event{}, fmt.Errorf("failed get event by id: %w", err)
	}

	return event, nil
}
