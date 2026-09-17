package booking

import (
	"bookingService/internal/seats"
	"context"
	"database/sql"
	"errors"
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
	inputData := ReserveInput{
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

func (s *BookingService) ConfirmBooking(ctx context.Context, bookingID int64, req ConfirmBookingRequest) (Result, error) {
	if bookingID <= 0 {
		return Result{}, ValidationError{Message: "invalid booking id"}
	}
	if req.UserID <= 0 {
		return Result{}, ValidationError{Message: "invalid user id"}
	}

	input := ConfirmInput{
		UserID:    req.UserID,
		BookingID: bookingID,
	}

	result, confirmErr := s.repo.ConfirmBooking(ctx, input)
	if confirmErr != nil {
		if errors.Is(confirmErr, sql.ErrNoRows) {
			return Result{}, ConflictError{Message: "failed update booking"}
		}

		log.Println(confirmErr)
		return Result{}, InternalError{Message: "failed to confirm booking"}
	}

	delReserveErr := s.reservationStore.DeleteReserve(ctx, result.SeatID)
	if delReserveErr != nil {
		log.Println(delReserveErr)
	}
	return result, nil
}

func (s *BookingService) CancelBooking(ctx context.Context, bookingID int64, req CancelBookingRequest) error {
	if bookingID <= 0 {
		return ValidationError{Message: "invalid booking id"}
	}
	if req.UserID <= 0 {
		return ValidationError{Message: "invalid user id"}
	}

	input := CancelInput{
		UserID:    req.UserID,
		BookingID: bookingID,
	}

	result, cancelErr := s.repo.CancelBooking(ctx, input)
	if cancelErr != nil {
		if errors.Is(cancelErr, sql.ErrNoRows) {
			return ConflictError{Message: "failed cancel booking"}
		}

		log.Println(cancelErr)
		return InternalError{Message: "failed to cancel booking"}
	}

	delReserveErr := s.reservationStore.DeleteReserve(ctx, result.SeatID)
	if delReserveErr != nil {
		log.Println(delReserveErr)
	}

	return nil
}

func validateInput(input ReserveInput) error {
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
