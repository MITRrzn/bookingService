package booking

import (
	"context"
	"time"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, input ReserveInput, reservedAt time.Time, ttl time.Time) (Result, error)
	ConfirmBooking(ctx context.Context, input ConfirmInput) (Result, error)
	CancelBooking(ctx context.Context, input CancelInput) (Result, error)
	GetBookingsByUser(ctx context.Context, userID int64) ([]ListItem, error)
	IsSeatConfirmed(ctx context.Context, seatID int64) (bool, error)
}
