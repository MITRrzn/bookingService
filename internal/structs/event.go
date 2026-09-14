package structs

import "time"

type Event struct {
	Name      string
	StartsAt  time.Time
	CreatedAt time.Time
}
