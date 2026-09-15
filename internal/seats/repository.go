package seats

import (
	"context"
)

type SeatRepository interface {
	AddSeatsToEvent(ctx context.Context, seats []SeatInput, eventID int64) error
	GetSeatsByEventID(ctx context.Context, eventID int64) ([]Seat, error)
}
