package seats

import "context"

type ServiceInterface interface {
	AddSeatsToEvent(ctx context.Context, seats []SeatInput, eventID int64) (int64, error)
	GetSeatsByEventID(ctx context.Context, id int64) ([]Seat, error)
}

type SeatRepository interface {
	AddSeats(ctx context.Context, seats []SeatInput, eventID int64) (int64, error)
	GetSeatsByEventID(ctx context.Context, eventID int64) ([]Seat, error)
	GetSeatByID(ctx context.Context, seatID int64) (Seat, error)
}
