package seats

import (
	"bookingService/internal/structs"
	"context"
)

type SeatsRepository interface {
	AddSeatsToEvent(ctx context.Context, seats []structs.SeatInput, eventId int64) error
	GetSeatsByEventId(ctx context.Context, eventId int64) ([]structs.Seat, error)
}
