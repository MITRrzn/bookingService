package booking

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockBookingService struct {
	result Result
	items  []ListItem
	error  error
}

func (m mockBookingService) ReserveSeat(ctx context.Context, eventID int64, seatId int64, req ReserveSeatsRequest) (Result, error) {
	return m.result, m.error
}

func (m mockBookingService) ConfirmBooking(ctx context.Context, bookingID int64, req ConfirmBookingRequest) (Result, error) {
	return m.result, m.error
}

func (m mockBookingService) CancelBooking(ctx context.Context, bookingID int64, req CancelBookingRequest) error {
	return m.error
}

func (m mockBookingService) ListBookings(ctx context.Context, userID int64) ([]ListItem, error) {
	return m.items, m.error
}

func TestReserveSeatHandlerErrors(t *testing.T) {
	testCases := []struct {
		name         string
		body         string
		seatID       any
		eventID      any
		error        error
		expectedCode int
	}{
		{
			name:         "invalid eventID",
			body:         `{"userId": 1}`,
			seatID:       1,
			eventID:      "qwe",
			error:        nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid seatID",
			body:         `{"userId": 1}`,
			seatID:       "rty",
			eventID:      1,
			error:        nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "incorrect json",
			body:         `{"userId": 1`,
			seatID:       2,
			eventID:      1,
			error:        nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "undefined seat",
			body:         `{"userId": 1}`,
			seatID:       2,
			eventID:      1,
			error:        sql.ErrNoRows,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "validation error, incorrect userID",
			body:         `{"userId": 0}`,
			seatID:       1,
			eventID:      1,
			error:        ValidationError{Message: "invalid user id"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "validation error, incorrect eventID",
			body:         `{"userId": 1}`,
			seatID:       1,
			eventID:      0,
			error:        ValidationError{Message: "invalid event id"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "validation error, incorrect seatID",
			body:         `{"userId": 1}`,
			seatID:       0,
			eventID:      1,
			error:        ValidationError{Message: "invalid seat id"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "conflict error",
			body:         `{"userId": 1}`,
			seatID:       1,
			eventID:      1,
			error:        ConflictError{Message: "failed update booking"},
			expectedCode: http.StatusConflict,
		},
		{
			name:         "internal server error",
			body:         `{"userId": 1}`,
			seatID:       1,
			eventID:      1,
			error:        InternalError{Message: "internal server error"},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &mockBookingService{error: tc.error}

			handler := NewHandler(service)
			mux := http.NewServeMux()

			mux.HandleFunc("POST /events/{eventID}/seats/{seatID}/reserve", handler.ReserveSeat)

			req := httptest.NewRequest(
				"POST",
				fmt.Sprintf("/events/%v/seats/%v/reserve", tc.eventID, tc.seatID),
				strings.NewReader(tc.body),
			)

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestReserveSeatHandlerSuccess(t *testing.T) {
	service := &mockBookingService{
		result: Result{
			BookingID:  1,
			UserID:     2,
			SeatID:     3,
			Status:     "reserved",
			ReservedAt: time.Time{},
			ExpiresAt:  time.Time{},
		},
	}

	handler := NewHandler(service)
	mux := http.NewServeMux()

	mux.HandleFunc("POST /events/{eventID}/seats/{seatID}/reserve", handler.ReserveSeat)

	req := httptest.NewRequest(
		"POST",
		fmt.Sprintf("/events/%d/seats/%d/reserve", 1, 2),
		strings.NewReader(`{"userId": 1}`),
	)

	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestConfirmHandlerErrors(t *testing.T) {
	testCases := []struct {
		name         string
		body         string
		bookingID    any
		error        error
		expectedCode int
	}{
		{
			name:         "invalid bookingID",
			body:         `{"userId": 1}`,
			bookingID:    "qwe",
			error:        nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid json",
			body:         `{"userId": 1`,
			bookingID:    1,
			error:        nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "validation error, incorrect bookingID",
			body:         `{"userId": 1}`,
			bookingID:    "0",
			error:        ValidationError{Message: "invalid booking id"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "validation error, incorrect userID",
			body:         `{"userId": 0}`,
			bookingID:    "1",
			error:        ValidationError{Message: "invalid user id"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "conflict error",
			body:         `{"userId": 0}`,
			bookingID:    "1",
			error:        ConflictError{Message: "failed update booking"},
			expectedCode: http.StatusConflict,
		},
		{
			name:         "internal server error",
			body:         `{"userId": 1}`,
			bookingID:    "1",
			error:        errors.New("internal server error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &mockBookingService{error: tc.error}

			handler := NewHandler(service)
			mux := http.NewServeMux()
			mux.HandleFunc("POST /bookings/{bookingID}/confirm", handler.Confirm)

			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("/bookings/%v/confirm", tc.bookingID),
				strings.NewReader(tc.body),
			)

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestConfirmHandlerSuccess(t *testing.T) {
	service := &mockBookingService{}

	handler := NewHandler(service)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /bookings/{bookingID}/confirm", handler.Confirm)

	req := httptest.NewRequest(
		http.MethodPost,
		"/bookings/1/confirm",
		strings.NewReader(`{"userId": 1}`),
	)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestCancelErrors(t *testing.T) {
	testCases := []struct {
		name         string
		body         string
		bookingID    any
		error        error
		expectedCode int
	}{
		{
			name:         "invalid booking id",
			body:         `{"userId": 1}`,
			bookingID:    "qwe",
			error:        nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid json",
			body:         `{"userId": 1`,
			bookingID:    "1",
			error:        nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "validation error, invalid booking id",
			body:         `{"userId": 1}`,
			bookingID:    "0",
			error:        ValidationError{Message: "invalid booking id"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "validation error, invalid user id",
			body:         `{"userId": 0}`,
			bookingID:    "1",
			error:        ValidationError{Message: "invalid user id"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "conflict error",
			body:         `{"userId": 0}`,
			bookingID:    "1",
			error:        ConflictError{Message: "failed cancel booking"},
			expectedCode: http.StatusConflict,
		},
		{
			name:         "internal server error",
			body:         `{"userId": 1}`,
			bookingID:    "1",
			error:        InternalError{Message: "failed cancel booking"},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &mockBookingService{
				error: tc.error,
			}

			handler := NewHandler(service)
			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /bookings/{bookingID}", handler.Cancel)

			req := httptest.NewRequest(
				http.MethodDelete,
				fmt.Sprintf("/bookings/%v", tc.bookingID),
				strings.NewReader(tc.body),
			)
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestCancelSuccess(t *testing.T) {
	service := &mockBookingService{}

	handler := NewHandler(service)
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /bookings/{bookingID}", handler.Cancel)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/bookings/1",
		strings.NewReader(`{"userId": 1}`),
	)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
}

func TestListHandlerErrors(t *testing.T) {
	testCases := []struct {
		name         string
		userId       any
		error        error
		expectedCode int
	}{
		{
			name:         "incorrect user id",
			userId:       "qwe",
			error:        nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "invalid user id",
			userId:       "0",
			error:        ValidationError{Message: "invalid user id"},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "internal server error",
			userId:       "1",
			error:        errors.New("internal server error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &mockBookingService{
				error: tc.error,
			}

			handler := NewHandler(service)
			mux := http.NewServeMux()
			mux.HandleFunc("GET /users/{userID}/bookings", handler.List)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%v/bookings", tc.userId), nil)
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestListHandlerEmpty(t *testing.T) {
	service := &mockBookingService{
		items: []ListItem{},
	}

	handler := NewHandler(service)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{userID}/bookings", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/users/1/bookings", nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestListHandlerSuccess(t *testing.T) {
	service := &mockBookingService{
		items: []ListItem{
			{
				BookingID:     1,
				BookingStatus: "reserved",
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
				SeatNumber:    "a2",
				SeatPrice:     110,
				EventName:     "test event",
				EventStartsAt: time.Time{},
			},
		},
	}

	handler := NewHandler(service)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{userID}/bookings", handler.List)

	req := httptest.NewRequest(http.MethodGet, "/users/1/bookings", nil)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}
