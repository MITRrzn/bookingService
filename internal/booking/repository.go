package booking

import (
	"context"
	"time"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, input Input, reservedAt time.Time, ttl time.Time) (Result, error)
	ConfirmBooking(ctx context.Context, bookingID int64) (Result, error)
}
