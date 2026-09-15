package seat

import (
	"bookingService/internal/structs"
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

func (r *Repository) AddSeatsToEvent(ctx context.Context, seats []structs.SeatInput, eventID int64) error {
	return nil
}

func (r *Repository) GetSeatsByEventID(ctx context.Context, eventID int64) ([]structs.Seat, error) {
	return nil, nil
}
