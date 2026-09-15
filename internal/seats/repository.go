package seats

import (
	"bookingService/internal/structs"
	"context"
)

type SeatRepository interface {
	AddSeatsToEvent(ctx context.Context, seats []structs.SeatInput, eventID int64) error
	GetSeatsByEventID(ctx context.Context, eventID int64) ([]structs.Seat, error)
}
