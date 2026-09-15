package event

import "time"

type Event struct {
	ID        int64
	Name      string
	StartsAt  time.Time
	CreatedAt time.Time
}
