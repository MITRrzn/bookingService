package seats

type AddSeatsRequest struct {
	Seats []SeatInput `json:"seats"`
}

type SeatInput struct {
	Number string
	Price  float64
}
