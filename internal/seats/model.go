package seats

import "time"

type Seat struct {
	ID        int64
	EventId   int64
	EventID   int64
	Number    string
	Price     float64
	CreatedAt time.Time
}
