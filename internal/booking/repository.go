package booking

import "context"

type BookingRepository interface {
	CreateBooking(ctx context.Context, input Input) (Result, error)
}
