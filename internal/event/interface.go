package event

import "context"

type ServiceInterface interface {
	CreateEvent(ctx context.Context, input CreateEventInput) (Event, error)
	GetActiveEvents(ctx context.Context) ([]Event, error)
	GetEventByID(ctx context.Context, id int64) (Event, error)
}

type EventRepository interface {
	CreateEvent(ctx context.Context, event CreateEventInput) (Event, error)
	GetActiveEvents(ctx context.Context) ([]Event, error)
	GetEventByID(ctx context.Context, id int64) (Event, error)
}
