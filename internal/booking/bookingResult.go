package booking

import "time"

type Result struct {
	BookingId  string    `json:"booking_id"`
	UserID     string    `json:"user_id"`
	SeatID     string    `json:"seat_id"`
	Status     string    `json:"status"`
	ReservedAt time.Time `json:"reserved_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}
