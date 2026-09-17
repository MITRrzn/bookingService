package booking

import (
	"context"
	"time"
)

type ReservationStore interface {
	Reserve(ctx context.Context, seatID int64, ttl time.Duration) (bool, error)

	DeleteReserve(ctx context.Context, seatID int64) error
}
