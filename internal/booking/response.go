package booking

import "time"

type Response struct {
	BookingId string    `json:"booking_id"`
	Status    string    `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
}
