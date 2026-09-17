package booking

type Response struct {
	Success     string `json:"success"`
	BookingData Result `json:"booking_data"`
}

type ListResponse struct {
	Success     string     `json:"success"`
	BookingData []ListItem `json:"booking_data"`
}
