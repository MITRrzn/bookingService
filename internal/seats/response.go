package seats

type SeatsResponse struct {
	Success string `json:"success"`
	Seats   []Seat `json:"seats"`
}
