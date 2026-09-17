package booking

import (
	"context"
	"time"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, input ReserveInput, reservedAt time.Time, ttl time.Time) (Result, error)
	ConfirmBooking(ctx context.Context, input ConfirmInput) (Result, error)
	CancelBooking(ctx context.Context, input CancelInput) (Result, error)
}
