package booking

type ReserveSeatsRequest struct {
	UserID int64 `json:"userId"`
}

type ConfirmBookingRequest struct {
	UserID int64 `json:"userId"`
}
