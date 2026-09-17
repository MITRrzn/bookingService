package booking

import (
	"bookingService/internal/seats"
	"context"
	"log"
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

func (s *BookingService) ReserveSeat(ctx context.Context, eventID int64, seatId int64, req ReserveSeatsRequest) (Result, error) {
	inputData := Input{
		UserID:  req.UserID,
		EventID: eventID,
		SeatID:  seatId,
	}

	validationErr := validateInput(inputData)
	if validationErr != nil {
		return Result{}, validationErr
	}

	seat, err := s.seatRepo.GetSeatByID(ctx, inputData.SeatID)
	if err != nil {
		return Result{}, err
	}
	if seat.EventID != inputData.EventID {
		return Result{}, NotFoundError{
			Message: "seat not found",
		}
	}

	now := time.Now()
	ttl := time.Minute * 5
	reserved, reserveErr := s.reservationStore.Reserve(ctx, inputData.SeatID, ttl)
	if reserveErr != nil {
		log.Println(reserveErr)
		return Result{}, InternalError{Message: "reservation service error"}
	}
	if !reserved {
		return Result{}, ConflictError{Message: "seat is already reserved"}
	}

	result, createErr := s.repo.CreateBooking(ctx, inputData, now, now.Add(ttl))
	if createErr != nil {
		log.Println(createErr)
		delReserveErr := s.reservationStore.DeleteReserve(ctx, inputData.SeatID)
		if delReserveErr != nil {
			log.Println(delReserveErr)
		}
		return Result{}, InternalError{Message: "create booking service error"}
	}

	return result, nil
}

func (s *BookingService) ConfirmBooking(ctx context.Context, bookingID int64) (Result, error) {
	if bookingID <= 0 {
		return Result{}, ValidationError{Message: "invalid booking id"}
	}

	result, confirmErr := s.repo.ConfirmBooking(ctx, bookingID)
	if confirmErr != nil {
		log.Println(confirmErr)
		return Result{}, InternalError{Message: "failed to confirm booking"}
	}

	delReserveErr := s.reservationStore.DeleteReserve(ctx, result.SeatID)
	if delReserveErr != nil {
		log.Println(delReserveErr)
	}
	return result, nil
}

func validateInput(input Input) error {
	if input.UserID <= 0 {
		log.Println("invalid user id", input)
		return ValidationError{Message: "invalid user id"}
	}
	if input.EventID <= 0 {
		log.Println("invalid event id", input)
		return ValidationError{Message: "invalid event id"}
	}
	if input.SeatID <= 0 {
		log.Println("invalid seat id", input)
		return ValidationError{Message: "invalid seat id"}
	}

	return nil
}
