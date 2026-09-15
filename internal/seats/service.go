package seats

type SeatService struct {
	repo SeatsRepository
}

func NewService(repo SeatsRepository) *SeatService {
	return &SeatService{
		repo: repo,
	}
}
