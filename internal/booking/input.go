package booking

type ReserveInput struct {
	UserID  int64
	EventID int64
	SeatID  int64
}

type ConfirmInput struct {
	UserID    int64
	BookingID int64
}

type CancelInput struct {
	UserID    int64
	BookingID int64
}
