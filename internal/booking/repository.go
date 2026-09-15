package booking

import "context"

type BookingRepository interface {
	ReserveSeat(ctx context.Context, input Input) error
}
