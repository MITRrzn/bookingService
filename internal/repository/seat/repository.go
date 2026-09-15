package seat

import (
	"bookingService/internal/structs"
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewSeatsRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) AddSeatsToEvent(ctx context.Context, seats []structs.SeatInput, eventId int64) error {
	return nil
}

func (r *Repository) GetSeatsByEventId(ctx context.Context, eventId int64) ([]structs.Seat, error) {
	return nil, nil
}
