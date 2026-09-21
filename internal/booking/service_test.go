package booking

import (
	"bookingService/internal/seats"
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockBookingRepository struct {
	booking          Result
	items            []ListItem
	isConfirmed      bool
	createBookingErr error
	isConfirmedErr   error
	confirmErr       error
	cancelErr        error
	listErr          error
}

type mockSeatRepository struct {
	seatsAmount int64
	seats       []seats.Seat
	seat        seats.Seat
	error       error
}

type mockReservationStore struct {
	reserved bool
	error    error
}

func (m mockSeatRepository) AddSeats(ctx context.Context, seats []seats.SeatInput, eventID int64) (int64, error) {
	return m.seatsAmount, m.error
}

func (m mockSeatRepository) GetSeatsByEventID(ctx context.Context, eventID int64) ([]seats.Seat, error) {
	return m.seats, m.error
}

func (m mockSeatRepository) GetSeatByID(ctx context.Context, seatID int64) (seats.Seat, error) {
	return m.seat, m.error
}

func (m mockReservationStore) Reserve(ctx context.Context, seatID int64, ttl time.Duration) (bool, error) {
	return m.reserved, m.error
}

func (m mockReservationStore) DeleteReserve(ctx context.Context, seatID int64) error {
	return m.error
}

func (m *mockBookingRepository) CreateBooking(ctx context.Context, input ReserveInput, reservedAt time.Time, ttl time.Time) (Result, error) {
	return m.booking, m.createBookingErr
}
func (m *mockBookingRepository) ConfirmBooking(ctx context.Context, input ConfirmInput) (Result, error) {
	return m.booking, m.confirmErr
}
func (m *mockBookingRepository) CancelBooking(ctx context.Context, input CancelInput) (Result, error) {
	return m.booking, m.cancelErr
}
func (m *mockBookingRepository) GetBookingsByUser(ctx context.Context, userID int64) ([]ListItem, error) {
	return m.items, m.listErr
}
func (m *mockBookingRepository) IsSeatConfirmed(ctx context.Context, seatID int64) (bool, error) {
	return m.isConfirmed, m.isConfirmedErr
}

func TestReserveSeatCreateBookingError(t *testing.T) {
	bookingRepo := &mockBookingRepository{
		isConfirmed:      false,
		createBookingErr: errors.New("create booking error"),
	}
	seatRepo := &mockSeatRepository{
		seat: seats.Seat{
			ID:        3,
			EventID:   2,
			Number:    "a1",
			Price:     100,
			CreatedAt: time.Time{},
		},
	}
	reservationStore := &mockReservationStore{
		reserved: true,
	}

	service := NewService(bookingRepo, seatRepo, reservationStore)
	_, err := service.ReserveSeat(context.Background(), 2, 3, ReserveSeatsRequest{UserID: 2})
	assert.Error(t, err)
	assert.Equal(t, InternalError{Message: "create booking service error"}, err)
}

func TestReserveSeatAlreadyReserved(t *testing.T) {
	bookingRepo := &mockBookingRepository{
		isConfirmed: false,
	}
	seatRepo := &mockSeatRepository{
		seat: seats.Seat{
			ID:        3,
			EventID:   2,
			Number:    "a1",
			Price:     100,
			CreatedAt: time.Time{},
		},
	}
	reservationStore := &mockReservationStore{}

	service := NewService(bookingRepo, seatRepo, reservationStore)
	_, err := service.ReserveSeat(context.Background(), 2, 3, ReserveSeatsRequest{UserID: 2})
	assert.Error(t, err)
	assert.Equal(t, ConflictError{Message: "seat is already reserved"}, err)
}

func TestReserveSeatReservationStoreError(t *testing.T) {
	bookingRepo := &mockBookingRepository{}
	seatRepo := &mockSeatRepository{
		seat: seats.Seat{
			ID:        3,
			EventID:   2,
			Number:    "a1",
			Price:     100,
			CreatedAt: time.Time{},
		},
	}
	reservationStore := &mockReservationStore{
		error: errors.New("redis error"),
	}

	service := NewService(bookingRepo, seatRepo, reservationStore)
	_, err := service.ReserveSeat(context.Background(), 2, 3, ReserveSeatsRequest{UserID: 2})
	assert.Error(t, err)
	assert.Equal(t, InternalError{Message: "reservation service error"}, err)
}

func TestReserveSeatAlreadyConfirmed(t *testing.T) {
	bookingRepo := &mockBookingRepository{
		isConfirmed: true,
	}
	seatRepo := &mockSeatRepository{
		seat: seats.Seat{
			ID:        3,
			EventID:   2,
			Number:    "a1",
			Price:     100,
			CreatedAt: time.Time{},
		},
	}
	reservationStore := &mockReservationStore{
		reserved: true,
	}

	service := NewService(bookingRepo, seatRepo, reservationStore)
	_, err := service.ReserveSeat(context.Background(), 2, 3, ReserveSeatsRequest{UserID: 2})
	assert.Error(t, err)
	assert.Equal(t, ConflictError{Message: "seat is already reserved"}, err)
}

func TestReserveSeatCheckConfirmedError(t *testing.T) {
	bookingRepo := &mockBookingRepository{
		isConfirmedErr: errors.New("isConfirmedErr"),
	}
	seatRepo := &mockSeatRepository{
		seat: seats.Seat{
			ID:        3,
			EventID:   2,
			Number:    "a1",
			Price:     100,
			CreatedAt: time.Time{},
		},
	}
	reservationStore := &mockReservationStore{}

	service := NewService(bookingRepo, seatRepo, reservationStore)
	_, err := service.ReserveSeat(context.Background(), 2, 3, ReserveSeatsRequest{UserID: 2})
	assert.Error(t, err)
	assert.Equal(t, InternalError{"failed check seat confirmation"}, err)
}

func TestReserveSeatSeatBelongsToAnotherEvent(t *testing.T) {
	bookingRepo := &mockBookingRepository{}
	seatRepo := &mockSeatRepository{
		seat: seats.Seat{
			ID:        1,
			EventID:   2,
			Number:    "a1",
			Price:     100,
			CreatedAt: time.Time{},
		},
	}
	reservationStore := &mockReservationStore{}

	service := NewService(bookingRepo, seatRepo, reservationStore)
	_, err := service.ReserveSeat(context.Background(), 1, 3, ReserveSeatsRequest{UserID: 2})
	assert.Error(t, err)
	assert.Equal(t, NotFoundError{Message: "seat not found"}, err)
}

func TestReserveSeatGetSeatError(t *testing.T) {
	expectedErr := errors.New("some error")

	bookingRepo := &mockBookingRepository{
		booking: Result{},
	}
	seatRepo := &mockSeatRepository{
		error: expectedErr,
	}
	reservationStore := &mockReservationStore{
		reserved: true,
	}

	service := NewService(bookingRepo, seatRepo, reservationStore)
	_, err := service.ReserveSeat(context.Background(), 1, 3, ReserveSeatsRequest{UserID: 2})
	assert.Equal(t, expectedErr, err)
}

func TestReserveSeatSuccess(t *testing.T) {
	expected := Result{
		BookingID:  1,
		UserID:     2,
		SeatID:     3,
		Status:     "reserved",
		ReservedAt: time.Time{},
		ExpiresAt:  time.Time{},
	}

	bookingRepo := &mockBookingRepository{
		booking: expected,
	}
	seatRepo := &mockSeatRepository{
		seat: seats.Seat{
			ID:        3,
			EventID:   1,
			Number:    "a1",
			Price:     100,
			CreatedAt: time.Time{},
		},
	}
	reservationStore := &mockReservationStore{
		reserved: true,
	}

	service := NewService(bookingRepo, seatRepo, reservationStore)

	res, err := service.ReserveSeat(context.Background(), 1, 3, ReserveSeatsRequest{UserID: 2})
	assert.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestListBookingsDBErr(t *testing.T) {
	bookingRepo := &mockBookingRepository{
		listErr: errors.New("db err"),
	}
	seatRepo := &mockSeatRepository{}
	reservationStore := &mockReservationStore{}

	service := NewService(bookingRepo, seatRepo, reservationStore)

	bookings, err := service.ListBookings(context.Background(), 1)
	assert.Nil(t, bookings)
	assert.Empty(t, bookings)
	assert.Error(t, err)
	assert.Equal(t, bookingRepo.listErr, err)
}

func TestConfirmBooking(t *testing.T) {
	cases := []struct {
		name           string
		bookingId      int64
		req            ConfirmBookingRequest
		wantErr        bool
		repoErr        error
		expectedErr    error
		expectedResult Result
	}{
		{
			name:        "confirm booking, success",
			bookingId:   1,
			req:         ConfirmBookingRequest{UserID: 1},
			wantErr:     false,
			repoErr:     nil,
			expectedErr: nil,
			expectedResult: Result{
				BookingID:  1,
				UserID:     1,
				SeatID:     2,
				Status:     "confirmed",
				ReservedAt: time.Time{},
				ExpiresAt:  time.Time{},
			},
		},
		{
			name:        "confirm booking, invalid booking id",
			bookingId:   0,
			req:         ConfirmBookingRequest{UserID: 1},
			wantErr:     true,
			repoErr:     nil,
			expectedErr: ValidationError{Message: "invalid booking id"},
		},
		{
			name:        "confirm booking, invalid user id",
			bookingId:   1,
			req:         ConfirmBookingRequest{UserID: 0},
			wantErr:     true,
			repoErr:     nil,
			expectedErr: ValidationError{Message: "invalid user id"},
		},
		{
			name:        "confirm booking, sql no rows",
			bookingId:   1,
			req:         ConfirmBookingRequest{UserID: 1},
			wantErr:     true,
			repoErr:     sql.ErrNoRows,
			expectedErr: ConflictError{Message: "failed update booking"},
		},
		{
			name:        "confirm booking, db err",
			bookingId:   1,
			req:         ConfirmBookingRequest{UserID: 1},
			wantErr:     true,
			repoErr:     errors.New("db err"),
			expectedErr: InternalError{Message: "failed to confirm booking"},
		},
	}

	for _, tc := range cases {
		bookingRepo := &mockBookingRepository{
			booking:    tc.expectedResult,
			confirmErr: tc.repoErr,
		}
		seatRepo := &mockSeatRepository{}
		reservationStore := &mockReservationStore{}
		service := NewService(bookingRepo, seatRepo, reservationStore)

		t.Run(tc.name, func(t *testing.T) {
			result, err := service.ConfirmBooking(context.Background(), tc.bookingId, tc.req)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			}
			if !tc.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestCancelBooking(t *testing.T) {
	cases := []struct {
		name        string
		bookingId   int64
		req         CancelBookingRequest
		wantErr     bool
		repoErr     error
		expectedErr error
	}{
		{
			name:      "cancel booking, success",
			bookingId: 1,
			req: CancelBookingRequest{
				UserID: 1,
			},
			wantErr:     false,
			repoErr:     nil,
			expectedErr: nil,
		},
		{
			name:      "cancel booking, invalid booking id",
			bookingId: 0,
			req: CancelBookingRequest{
				UserID: 1,
			},
			wantErr:     true,
			repoErr:     nil,
			expectedErr: ValidationError{Message: "invalid booking id"},
		},
		{
			name:      "cancel booking, invalid user id",
			bookingId: 1,
			req: CancelBookingRequest{
				UserID: 0,
			},
			wantErr:     true,
			repoErr:     nil,
			expectedErr: ValidationError{Message: "invalid user id"},
		},
		{
			name:      "cancel booking, no rows err",
			bookingId: 1,
			req: CancelBookingRequest{
				UserID: 1,
			},
			wantErr:     true,
			repoErr:     sql.ErrNoRows,
			expectedErr: ConflictError{Message: "failed cancel booking"},
		},
		{
			name:      "cancel booking, db error",
			bookingId: 1,
			req: CancelBookingRequest{
				UserID: 1,
			},
			wantErr:     true,
			repoErr:     errors.New("db err"),
			expectedErr: InternalError{Message: "failed to cancel booking"},
		},
	}

	for _, tc := range cases {
		bookingRepo := &mockBookingRepository{
			booking: Result{
				SeatID: 10,
			},
			cancelErr: tc.repoErr,
		}
		seatRepo := &mockSeatRepository{}
		reservationStore := &mockReservationStore{}
		service := NewService(bookingRepo, seatRepo, reservationStore)

		t.Run(tc.name, func(t *testing.T) {
			err := service.CancelBooking(context.Background(), tc.bookingId, tc.req)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			}
			if !tc.wantErr {
				assert.NoError(t, err)
			}
		})
	}
}

func TestListBookingsEmptyResult(t *testing.T) {
	bookingRepo := &mockBookingRepository{
		items: []ListItem{},
	}
	seatRepo := &mockSeatRepository{}
	reservationStore := &mockReservationStore{}

	service := NewService(bookingRepo, seatRepo, reservationStore)

	bookings, err := service.ListBookings(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, bookings)
	assert.Empty(t, bookings)
}

func TestListBookingsSuccess(t *testing.T) {
	expected := []ListItem{
		{
			BookingID:     1,
			BookingStatus: "confirmed",
			ReservedAt:    time.Time{},
			ExpiresAt:     time.Time{},
			ConfirmedAt:   nil,
			SeatNumber:    "a1",
			SeatPrice:     100,
			EventName:     "test event",
			EventStartsAt: time.Time{},
		},
		{
			BookingID:     2,
			BookingStatus: "confirmed",
			ReservedAt:    time.Time{},
			ExpiresAt:     time.Time{},
			ConfirmedAt:   nil,
			SeatNumber:    "b1",
			SeatPrice:     150,
			EventName:     "another event",
			EventStartsAt: time.Time{},
		},
	}
	bookingRepo := &mockBookingRepository{
		items: expected,
	}
	seatRepo := &mockSeatRepository{}
	reservationStore := &mockReservationStore{}

	service := NewService(bookingRepo, seatRepo, reservationStore)

	bookings, err := service.ListBookings(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, bookings)
}

func TestValidateInput(t *testing.T) {
	cases := []struct {
		name        string
		input       ReserveInput
		wantErr     bool
		expectedErr error
	}{
		{
			name: "validate success",
			input: ReserveInput{
				UserID:  1,
				EventID: 2,
				SeatID:  3,
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			name: "invalid user ID",
			input: ReserveInput{
				UserID:  0,
				EventID: 1,
				SeatID:  2,
			},
			wantErr:     true,
			expectedErr: ValidationError{Message: "invalid user id"},
		},
		{
			name: "invalid event ID",
			input: ReserveInput{
				UserID:  1,
				EventID: 0,
				SeatID:  2,
			},
			wantErr:     true,
			expectedErr: ValidationError{Message: "invalid event id"},
		},
		{
			name: "invalid seat ID",
			input: ReserveInput{
				UserID:  1,
				EventID: 2,
				SeatID:  0,
			},
			wantErr:     true,
			expectedErr: ValidationError{Message: "invalid seat id"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateInput(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			}
			if !tc.wantErr {
				assert.NoError(t, err)
				assert.Nil(t, err)
			}
		})
	}
}
