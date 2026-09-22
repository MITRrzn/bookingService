package booking

import (
	"context"
	"time"
)

type ServiceInterface interface {
	ReserveSeat(ctx context.Context, eventID int64, seatId int64, req ReserveSeatsRequest) (Result, error)
	ConfirmBooking(ctx context.Context, bookingID int64, req ConfirmBookingRequest) (Result, error)
	CancelBooking(ctx context.Context, bookingID int64, req CancelBookingRequest) error
	ListBookings(ctx context.Context, userID int64) ([]ListItem, error)
}

type BookingRepository interface {
	CreateBooking(ctx context.Context, input ReserveInput, reservedAt time.Time, ttl time.Time) (Result, error)
	ConfirmBooking(ctx context.Context, input ConfirmInput) (Result, error)
	CancelBooking(ctx context.Context, input CancelInput) (Result, error)
	GetBookingsByUser(ctx context.Context, userID int64) ([]ListItem, error)
	IsSeatConfirmed(ctx context.Context, seatID int64) (bool, error)
}

type ReservationStore interface {
	Reserve(ctx context.Context, seatID int64, ttl time.Duration) (bool, error)
	DeleteReserve(ctx context.Context, seatID int64) error
}
