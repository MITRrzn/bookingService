package booking

import "time"

type Result struct {
	BookingID  int64     `json:"booking_id"`
	UserID     int64     `json:"user_id"`
	SeatID     int64     `json:"seat_id"`
	Status     string    `json:"status"`
	ReservedAt time.Time `json:"reserved_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}
