package seats

type GetSeatsResponse struct {
	Success string `json:"success"`
	Seats   []Seat `json:"seats"`
}

type PostSeatsResponse struct {
	Success string `json:"success"`
	Amount  int64  `json:"amount"`
}
