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

type ListItem struct {
	BookingID     int64      `json:"booking_id"`
	BookingStatus string     `json:"booking_status"`
	ReservedAt    time.Time  `json:"reserved_at"`
	ExpiresAt     time.Time  `json:"expires_at"`
	ConfirmedAt   *time.Time `json:"confirmed_at"`
	SeatNumber    string     `json:"seat_number"`
	SeatPrice     float64    `json:"seat_price"`
	EventName     string     `json:"event_name"`
	EventStartsAt time.Time  `json:"event_starts_at"`
}
