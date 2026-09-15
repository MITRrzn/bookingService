package seat

import (
	"bookingService/internal/seats"
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewSeatRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) AddSeatsToEvent(ctx context.Context, seats []seats.SeatInput, eventID int64) error {
	return nil
}

func (r *Repository) GetSeatsByEventID(ctx context.Context, eventID int64) ([]seats.Seat, error) {
	rows, queryErr := r.db.QueryContext(
		ctx,
		`
		SELECT * FROM seats
		WHERE id = $1
		`,
		eventID,
	)

	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	var seatsList []seats.Seat
	for rows.Next() {
		var model seats.Seat
		scanErr := rows.Scan(&model)
		if scanErr != nil {
			return nil, scanErr
		}

		seatsList = append(seatsList, model)
	}

	return seatsList, nil
}
