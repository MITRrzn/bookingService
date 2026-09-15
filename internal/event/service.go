package event

import (
	"bookingService/internal/structs"
	"context"
	"strings"
	"unicode/utf8"
)

type EventService struct {
	repo EventRepository
}

func NewService(repo EventRepository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) CreateEvent(ctx context.Context, input CreateEventInput) (Event, error) {
	validationErr := validateCreateEventInput(input)
	if validationErr != nil {
		return Event{}, validationErr
	}
	input.Name = strings.TrimSpace(input.Name)

	return s.repo.CreateEvent(ctx, input)
}

func validateCreateEventInput(input CreateEventInput) error {
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

func (s *EventService) GetEvents(ctx context.Context) ([]Event, error) {
	events, err := s.repo.GetActiveEvents(ctx)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *EventService) GetEventByID(ctx context.Context, id int64) (Event, error) {
	return s.repo.GetEventByID(ctx, id)
}
