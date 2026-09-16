package booking

import (
	"bookingService/internal/booking"
	"context"
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func NewBookingRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) CreateBooking(ctx context.Context, input booking.Input) (booking.Result, error) {
	var result booking.Result
	now := time.Now()

	err := r.db.QueryRowContext(
		ctx,
		`
        INSERT INTO bookings (
            user_id,
            seat_id,
            status,
            reserved_at,
            expires_at
        )
        VALUES ($1, $2, $3, $4, $5)
        RETURNING
            id,
            user_id,
            seat_id,
            status,
            reserved_at,
            expires_at
    	`,
		input.UserID,
		input.SeatID,
		"reserved",
		now,
		now.Add(time.Minute*5),
	).Scan(
		&result.BookingID,
		&result.UserID,
		&result.SeatID,
		&result.Status,
		&result.ReservedAt,
		&result.ExpiresAt,
	)
	if err != nil {
		return booking.Result{}, err
	}

	return result, nil
}
