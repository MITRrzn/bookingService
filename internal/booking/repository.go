package booking

import (
	"context"
	"time"
)

type BookingRepository interface {
	CreateBooking(ctx context.Context, input Input, reservedAt time.Time, ttl time.Time) (Result, error)
}
