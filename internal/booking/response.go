package booking

type Response struct {
	Success     string `json:"success"`
	BookingData Result `json:"booking_data"`
}
