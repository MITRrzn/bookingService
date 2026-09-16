package booking

import (
	"context"
	"time"
)

type ReservationStore interface {
	Reserve(
		ctx context.Context,
		input Input,
		ttl time.Duration,
	) (bool, error)
}
