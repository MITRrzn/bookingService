package booking

import (
	"context"
	"fmt"
)

type BookingService struct {
	repo BookingRepository
}

func NewService(repo BookingRepository) *BookingService {
	return &BookingService{
		repo: repo,
	}
}

func ReserveSeat(ctx context.Context, req ReserveSeatsRequest, eventID int64, seatID string) error {
	inputData := ReserveSeatInput{
		UserId:  req.UserId,
		SeatId:  seatID,
		EventID: eventID,
	}

	validErr := validateReserveSeats(inputData)
	if validErr != nil {
		return fmt.Errorf("invalid request data: %w", validErr)
	}

	return nil
}

func validateReserveSeats(input ReserveSeatInput) error {
	if input.UserId <= 0 {
		return ValidationError{
			Message: "userId must be greater than zero",
		}
	}

	if input.EventID <= 0 {
		return ValidationError{
			Message: "eventID must be greater than zero",
		}
	}

	if input.SeatId == "" {
		return ValidationError{
			Message: "seatID is required",
		}
	}

	return nil
}
