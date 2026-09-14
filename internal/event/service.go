package event

import (
	"bookingService/internal/structs"
	"context"
	"strings"
	"unicode/utf8"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateEvent(ctx context.Context, input structs.CreateEventInput) (structs.Event, error) {
	validationErr := validateCreateEventInput(input)
	if validationErr != nil {
		return structs.Event{}, validationErr
	}
	input.Name = strings.TrimSpace(input.Name)

	return s.repo.CreateEvent(ctx, input)
}

func validateCreateEventInput(input structs.CreateEventInput) error {
	name := strings.TrimSpace(input.Name)

	if name == "" {
		return structs.ValidationError{
			Message: "empty event name",
		}
	}

	if utf8.RuneCountInString(name) < 3 {
		return structs.ValidationError{
			Message: "event name is too short",
		}
	}

	if utf8.RuneCountInString(name) > 255 {
		return structs.ValidationError{
			Message: "event name is too long",
		}
	}

	return nil
}

//func (s *Service) GetEvents(ctx context.Context) ([]Event, error) {
//}

func (s *Service) GetEventByID(ctx context.Context, id int64) (structs.Event, error) {
	return s.repo.GetEventByID(ctx, id)
}
