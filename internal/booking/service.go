package booking

import (
	"bookingService/internal/seats"
	"context"
	"time"
)

type BookingService struct {
	repo             BookingRepository
	seatRepo         seats.SeatRepository
	reservationStore ReservationStore
}

func NewService(repo BookingRepository, seatRepo seats.SeatRepository, reservationStore ReservationStore) *BookingService {
	return &BookingService{
		repo:             repo,
		seatRepo:         seatRepo,
		reservationStore: reservationStore,
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

	seat, err := s.seatRepo.GetSeatByID(ctx, inputData.SeatID)
	if err != nil {
		return err
	}
	if seat.EventID != inputData.EventID {
		return NotFoundError{
			Message: "seat not found",
		}
	}

	reserved, reserveErr := s.reservationStore.Reserve(ctx, inputData, time.Minute*5)
	if reserveErr != nil {
		return InternalError{Message: "reservation service error"}
	}
	if !reserved {
		return ConflictError{Message: "seat is already reserved"}
	}

	//Todo store to bookings table

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
