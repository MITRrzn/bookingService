package seats

type AddSeatsRequest struct {
	Seats []SeatInput `json:"seats"`
}

type SeatInput struct {
	Number string  `json:"number"`
	Price  float64 `json:"price"`
}
