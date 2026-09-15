package seats

import (
	"context"
)

type SeatRepository interface {
	AddSeats(ctx context.Context, seats []SeatInput, eventID int64) (int64, error)
	GetSeatsByEventID(ctx context.Context, eventID int64) ([]Seat, error)
}
