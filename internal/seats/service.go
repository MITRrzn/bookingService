package seats

type SeatService struct {
	repo SeatRepository
}

func NewService(repo SeatRepository) *SeatService {
	return &SeatService{
		repo: repo,
	}
}
