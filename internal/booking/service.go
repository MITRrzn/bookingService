package booking

import (
	"bookingService/internal/seats"
	"context"
)

type BookingService struct {
	repo     BookingRepository
	seatRepo seats.SeatRepository
}

func NewService(repo BookingRepository) *BookingService {
	return &BookingService{
		repo: repo,
	}
}

func (s *BookingService) ReserveSeat(ctx context.Context, eventID int64, seatId int64, req ReserveSeatsRequest) error {
	inputData := Input{
		UserID:  req.UserID,
		EventID: eventID,
		SeatID:  seatId,
	}

	validationErr := validateInput(inputData)
	if validationErr != nil {
		return validationErr
	}

	_, err := s.seatRepo.GetSeatByID(ctx, inputData.SeatID)
	if err != nil {
		return err
	}

	return nil
}

func validateInput(input Input) error {
	if input.UserID <= 0 {
		return ValidationError{Message: "invalid user id"}
	}
	if input.EventID <= 0 {
		return ValidationError{Message: "invalid event id"}
	}
	if input.SeatID <= 0 {
		return ValidationError{Message: "invalid seat id"}
	}

	return nil
}
