package seats

import (
	"context"
	"fmt"
	"strings"
)

type SeatService struct {
	repo SeatRepository
}

func NewService(repo SeatRepository) *SeatService {
	return &SeatService{
		repo: repo,
	}
}

func (s *SeatService) AddSeatsToEvent(ctx context.Context, seats []SeatInput, eventID int64) (int64, error) {
	if eventID <= 0 {
		return 0, ValidationError{
			Message: "eventID must be greater than zero",
		}
	}

	seatsValidation := validateSeatsList(seats)
	if seatsValidation == nil {
		return 0, fmt.Errorf("invalid seats list: %v", seatsValidation)
	}

	return s.repo.AddSeats(ctx, seats, eventID)
}

func (s *SeatService) GetSeatsByEventID(ctx context.Context, id int64) ([]Seat, error) {
	if id <= 0 {
		return nil, ValidationError{
			Message: "eventID must be greater than zero",
		}
	}

	seats, err := s.repo.GetSeatsByEventID(ctx, id)
	if err != nil {
		return nil, err
	}

	return seats, nil
}

func validateSeatsList(seats []SeatInput) error {
	if len(seats) == 0 {
		return ValidationError{
			Message: "no seats to add",
		}
	}

	for _, seat := range seats {
		if seat.Price <= 0 {
			return ValidationError{
				Message: "invalid price",
			}
		}
	}

	seen := make(map[string]struct{})

	for _, seat := range seats {
		number := strings.TrimSpace(seat.Number)

		if _, exists := seen[number]; exists {
			return ValidationError{
				Message: fmt.Sprintf("duplicate seat number: %v", number),
			}
		}

		seen[number] = struct{}{}
	}

	return nil
}
