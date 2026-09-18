package seats

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockSeatRepository struct {
	amount int64
	seats  []Seat
	seat   Seat
	err    error
}

func (m *mockSeatRepository) AddSeats(
	ctx context.Context,
	seats []SeatInput,
	eventID int64,
) (int64, error) {
	return m.amount, m.err
}

func (m *mockSeatRepository) GetSeatsByEventID(
	ctx context.Context,
	eventID int64,
) ([]Seat, error) {
	return m.seats, m.err
}

func (m *mockSeatRepository) GetSeatByID(
	ctx context.Context,
	seatID int64,
) (Seat, error) {
	return m.seat, m.err
}

func TestAddSeatsToEventSuccess(t *testing.T) {
	repo := &mockSeatRepository{
		amount: 2,
		err:    nil,
	}
	service := NewService(repo)
	seats := []SeatInput{
		{Number: "A1", Price: 100},
		{Number: "A2", Price: 150},
	}

	amount, err := service.AddSeatsToEvent(context.Background(), seats, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), amount)
}

func TestAddSeatsToEventInvalidEvenID(t *testing.T) {
	repo := &mockSeatRepository{}
	service := NewService(repo)

	_, err := service.AddSeatsToEvent(context.Background(), []SeatInput{}, -1)
	assert.Equal(t, ValidationError{"eventID must be greater than zero"}, err)
}

func TestGetSeatsByEventIDInvalidID(t *testing.T) {
	repo := &mockSeatRepository{}
	service := NewService(repo)

	_, err := service.GetSeatsByEventID(context.Background(), -1)
	if err == nil {
		t.Fatal("expected error")
	}
	assert.Equal(t, ValidationError{"eventID must be greater than zero"}.Error(), err.Error())
}

func TestGetSeatsByEventIDSuccess(t *testing.T) {
	repo := &mockSeatRepository{
		seats: []Seat{
			{
				ID:      1,
				EventID: 10,
				Number:  "A1",
				Price:   100,
			},
			{
				ID:      2,
				EventID: 10,
				Number:  "A2",
				Price:   150,
			},
		},
		err: nil,
	}

	service := NewService(repo)

	result, err := service.GetSeatsByEventID(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 seats, got %d", len(result))
	}

	if result[0].Number != "A1" {
		t.Errorf("expected first seat A1, got %s", result[0].Number)
	}

	if result[1].Number != "A2" {
		t.Errorf("expected second seat A2, got %s", result[1].Number)
	}
}

func TestValidateSeatsList(t *testing.T) {
	cases := []struct {
		name    string
		seats   []SeatInput
		wantErr bool
	}{
		{
			name: "success",
			seats: []SeatInput{
				{Number: "A1", Price: 100},
				{Number: "A2", Price: 150},
			},
			wantErr: false,
		},
		{
			name:    "empty seats",
			seats:   []SeatInput{},
			wantErr: true,
		},
		{
			name: "invalid price",
			seats: []SeatInput{
				{Number: "A1", Price: 0},
				{Number: "A2", Price: 150},
			},
			wantErr: true,
		},
		{
			name: "empty seat number",
			seats: []SeatInput{
				{Number: "A1", Price: 100},
				{Number: "", Price: 150},
			},
			wantErr: true,
		},
		{
			name: "duplicate seat",
			seats: []SeatInput{
				{Number: "A1", Price: 100},
				{Number: "A1", Price: 150},
			},
			wantErr: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			err := validateSeatsList(testCase.seats)
			if testCase.wantErr && err == nil {
				t.Error("expected error but got nil")
			}

			if !testCase.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
