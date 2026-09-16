package ReservationStore

import (
	"context"
	"time"
)

type ReservationStore interface {
	Reserve(
		ctx context.Context,
		eventID int64,
		seatID int64,
		userID int64,
		ttl time.Duration,
	) (bool, error)
}
