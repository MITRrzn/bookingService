package seats

import (
	"context"
)

type SeatService struct {
	repo SeatRepository
}

func NewService(repo SeatRepository) *SeatService {
	return &SeatService{
		repo: repo,
	}
}

func (s *SeatService) GetSeatsByEventID(ctx context.Context, id int64) ([]Seat, error) {
	seats, err := s.repo.GetSeatsByEventID(ctx, id)
	if err != nil {
		return nil, err
	}

	return seats, nil
}
