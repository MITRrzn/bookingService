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

func (r Repository) CreateBooking(ctx context.Context, input booking.ReserveInput, reservedAt time.Time, ttl time.Time) (booking.Result, error) {
	var result booking.Result

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
		reservedAt,
		ttl,
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

func (r Repository) ConfirmBooking(ctx context.Context, input booking.ConfirmInput) (booking.Result, error) {
	var result booking.Result

	err := r.db.QueryRowContext(
		ctx,
		`
			UPDATE bookings SET status = $1, confirmed_at = NOW()
			WHERE id = $2 AND status = $3 AND expires_at > NOW() AND user_id = $4
			RETURNING
    		id,
    		user_id,
    		seat_id,
    		status,
    		reserved_at,
    		expires_at
    	`,
		"confirmed",
		input.BookingID,
		"reserved",
		input.UserID,
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
