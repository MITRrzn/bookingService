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
		SELECT id, event_id, number, price, created_at FROM seats
		WHERE event_id = $1
		`,
		eventID,
	)

	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	seatsList := make([]seats.Seat, 0)
	for rows.Next() {
		var model seats.Seat
		scanErr := rows.Scan(&model.ID, &model.EventID, &model.Number, &model.Price, &model.CreatedAt)
		if scanErr != nil {
			return nil, scanErr
		}

		seatsList = append(seatsList, model)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return seatsList, nil
}
