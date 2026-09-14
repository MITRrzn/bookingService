package event

import (
	"bookingService/internal/structs"
	"context"
)

type Service struct {
	repo Repository
}

func (s *Service) CreateEvent(ctx context.Context, input structs.CreateEventInput) (structs.Event, error) {
	return s.repo.CreateEvent(ctx, input)
}

//func (s *Service) GetEvents(ctx context.Context) ([]Event, error) {
//}

//func (s *Service) GetEventByID(ctx context.Context, id int64) (Event, error) {
//}
