package booking

import "context"

type BookingService struct {
	repo BookingRepository
}

func NewService(repo BookingRepository) *BookingService {
	return &BookingService{
		repo: repo,
	}
}

func (s *BookingService) AddSeatsToEvent(ctx context.Context, eventID int64, seatId string, req ReserveSeatsRequest) error {
	inputData := Input{
		UserID:  req.UserID,
		EventID: eventID,
		SeatID:  seatId,
	}

	validationErr := validateInput(inputData)
	if validationErr != nil {
		return validationErr
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
	if input.SeatID == "" {
		return ValidationError{Message: "invalid seat id"}
	}

	return nil
}
