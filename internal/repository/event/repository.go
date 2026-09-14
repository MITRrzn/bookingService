package event

import (
	"bookingService/internal/structs"
	"context"
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
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
		return structs.Event{}, errors.New("error creating event")
	}

	return event, nil
}
