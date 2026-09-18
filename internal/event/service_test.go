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

func TestCreateEventSuccess(t *testing.T) {
	expectedEvent := Event{
		ID:   1,
		Name: "success event",
	}

	repo := &mockEventRepository{
		event: expectedEvent,
	}
	service := NewService(repo)

	result, err := service.CreateEvent(context.Background(), CreateEventInput{
		Name:     "success event",
		StartsAt: "2026-11-12 11:12:13",
	})

	assert.NoError(t, err)
	assert.Equal(t, expectedEvent, result)
}

func TestValidateCreateEventInput(t *testing.T) {
	cases := []struct {
		name        string
		input       CreateEventInput
		wantErr     bool
		expectedErr error
	}{
		{
			name: "valid input",
			input: CreateEventInput{
				Name:     "test input",
				StartsAt: "2026-11-12 11:12:13",
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "empty input event name",
			input: CreateEventInput{
				Name:     "",
				StartsAt: "2026-11-12 11:12:13",
			},
			wantErr:     true,
			expectedErr: ValidationError{Message: "empty event name"},
		},
		{
			name: "short input event name",
			input: CreateEventInput{
				Name:     "ti",
				StartsAt: "2026-11-12 11:12:13",
			},
			wantErr:     true,
			expectedErr: ValidationError{Message: "event name is too short"},
		},
		{
			name: "long input event name",
			input: CreateEventInput{
				Name:     "qwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnmqwertyuiopasdfghjklzxcvbnm",
				StartsAt: "2026-11-12 11:12:13",
			},
			wantErr:     true,
			expectedErr: ValidationError{Message: "event name is too long"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateEventInput(tt.input)
			if tt.wantErr {
				assert.Equal(t, tt.expectedErr, err)
			}
			if !tt.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetActiveEventsEmptyEvents(t *testing.T) {
	repo := &mockEventRepository{
		listEvents: []Event{},
	}

	service := NewService(repo)
	res, err := service.GetActiveEvents(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Empty(t, res)
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
