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

func (r Repository) CancelBooking(ctx context.Context, input booking.CancelInput) (booking.Result, error) {
	var result booking.Result

	err := r.db.QueryRowContext(
		ctx,
		`
			UPDATE bookings SET status = $1
			WHERE id = $2 AND status = $3 AND expires_at > NOW() AND user_id = $4
			RETURNING
    		id,
    		user_id,
    		seat_id,
    		status,
    		reserved_at,
    		expires_at
    	`,
		"cancelled",
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

func (r Repository) GetBookingsByUser(ctx context.Context, userID int64) ([]booking.ListItem, error) {
	rows, err := r.db.QueryContext(ctx,
		`
			SELECT
			b.id,
			b.status,
			b.reserved_at,
			b.expires_at,
			b.confirmed_at,
			s.number,
			s.price,
			e.name,
			e.starts_at
			FROM bookings b JOIN seats s ON b.seat_id = s.id JOIN events e ON s.event_id = e.id
				WHERE b.user_id = $1
			ORDER BY b.reserved_at DESC `,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	items := make([]booking.ListItem, 0)
	for rows.Next() {
		var item booking.ListItem
		scanErr := rows.Scan(
			&item.BookingID,
			&item.BookingStatus,
			&item.ReservedAt,
			&item.ExpiresAt,
			&item.ConfirmedAt,
			&item.SeatNumber,
			&item.SeatPrice,
			&item.EventName,
			&item.EventStartsAt,
		)
		if scanErr != nil {
			return nil, scanErr
		}

		items = append(items, item)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return items, nil
}

func (r Repository) UpdateExpiredBookings(ctx context.Context) (int64, error) {
	result, execErr := r.db.ExecContext(ctx, `UPDATE bookings SET status = 'expired' WHERE status = 'reserved' AND expires_at <= now()`)
	if execErr != nil {
		return 0, execErr
	}

	amount, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return amount, nil
}
