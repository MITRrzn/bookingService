package event

import (
	"context"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event CreateEventInput) (Event, error)
	GetActiveEvents(ctx context.Context) ([]Event, error)
	GetEventByID(ctx context.Context, id int64) (Event, error)
}
