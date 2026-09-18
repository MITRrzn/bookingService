package event

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockEventRepository struct {
	event      Event
	listEvents []Event
	err        error
}

func (m *mockEventRepository) CreateEvent(ctx context.Context, event CreateEventInput) (Event, error) {
	return m.event, m.err
}

func (m *mockEventRepository) GetActiveEvents(ctx context.Context) ([]Event, error) {
	return m.listEvents, m.err
}

func (m *mockEventRepository) GetEventByID(ctx context.Context, id int64) (Event, error) {
	return m.event, m.err
}

func TestGetActiveEventsNilEvents(t *testing.T) {
	repo := &mockEventRepository{
		listEvents: nil,
	}

	service := NewService(repo)
	res, err := service.GetActiveEvents(context.Background())
	assert.NoError(t, err)
	assert.Nil(t, res)
}

func TestGetActiveEventsDbError(t *testing.T) {
	repo := &mockEventRepository{
		err: errors.New("GetActiveEventsDbError"),
	}
	service := NewService(repo)

	_, err := service.GetActiveEvents(context.Background())
	assert.Error(t, err)
}

func TestGetActiveEventsSuccess(t *testing.T) {
	expectedEvents := []Event{
		{
			ID:   1,
			Name: "test event 1",
		},
		{
			ID:   2,
			Name: "test event 2",
		},
	}

	repo := &mockEventRepository{
		listEvents: expectedEvents,
	}

	service := NewService(repo)

	result, err := service.GetActiveEvents(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedEvents, result)
}

func TestGetEventByIDSuccess(t *testing.T) {
	expected := Event{
		ID:   1,
		Name: "test event",
	}

	repo := &mockEventRepository{
		event: expected,
	}

	service := NewService(repo)

	result, err := service.GetEventByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetEventByIDInvalidID(t *testing.T) {
	repo := &mockEventRepository{}

	service := NewService(repo)

	_, err := service.GetEventByID(context.Background(), -1)
	assert.Error(t, err)
	assert.Equal(t, ValidationError{"invalid event id"}, err)
}

func TestGetEventByIDRepositoryError(t *testing.T) {
	repo := &mockEventRepository{
		err: errors.New("database error"),
	}

	service := NewService(repo)

	_, err := service.GetEventByID(context.Background(), 2)
	assert.Error(t, err)
	assert.Equal(t, repo.err, err)
}
