package event

import (
	"bookingService/internal/structs"
	"context"
)

type Repository interface {
	CreateEvent(ctx context.Context, event structs.CreateEventInput) (structs.Event, error)
	GetActiveEvents(ctx context.Context) ([]structs.Event, error)
	GetEventByID(ctx context.Context, id int64) (structs.Event, error)
}
