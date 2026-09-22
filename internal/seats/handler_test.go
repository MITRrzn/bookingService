package seats

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockSeatService struct {
	seats  []Seat
	seat   Seat
	amount int64
	err    error
}

func (m MockSeatService) AddSeatsToEvent(ctx context.Context, seats []SeatInput, eventID int64) (int64, error) {
	return m.amount, m.err
}

func (m MockSeatService) GetSeatsByEventID(ctx context.Context, id int64) ([]Seat, error) {
	return m.seats, m.err
}

func TestAddSeatsToEventHandlerErrors(t *testing.T) {
	cases := []struct {
		name         string
		id           string
		body         string
		expectedCode int
		error        error
	}{
		{
			name:         "invalid eventID",
			id:           "qwe",
			body:         `{"seats": [{"number": "a1","price": 10},{"number": "a2","price": 20}]}`,
			expectedCode: http.StatusBadRequest,
			error:        nil,
		},
		{
			name:         "invalid json",
			id:           "1",
			body:         `{"seats": `,
			expectedCode: http.StatusBadRequest,
			error:        nil,
		},
		{
			name:         "validation error, no seats",
			id:           "1",
			body:         `{"seats": []}`,
			expectedCode: http.StatusBadRequest,
			error:        ValidationError{Message: "no seats to add"},
		},
		{
			name:         "validation error, invalid price",
			id:           "1",
			body:         `{"seats": [{"number": "a1","price": -1}]}`,
			expectedCode: http.StatusBadRequest,
			error:        ValidationError{Message: "invalid price"},
		},
		{
			name:         "validation error, empty seat number",
			id:           "1",
			body:         `{"seats": [{"number": "","price": 10}]}`,
			expectedCode: http.StatusBadRequest,
			error:        ValidationError{Message: "seat number is empty"},
		},
		{
			name:         "validation error, duplicate seat",
			id:           "1",
			body:         `{"seats": [{"number": "a1","price": 10},{"number": "a2","price": 20},{"number": "a2","price": 30}]}`,
			expectedCode: http.StatusBadRequest,
			error:        ValidationError{Message: "duplicate seat number: a2"},
		},
		{
			name:         "internal server error",
			id:           "1",
			body:         `{"seats": [{"number": "","price": 10}]}`,
			expectedCode: http.StatusInternalServerError,
			error:        errors.New("internal server error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := &MockSeatService{err: tc.error}
			handler := NewHandler(service)
			mux := http.NewServeMux()
			mux.HandleFunc("POST /events/{eventID}/seats", handler.AddSeatsToEvent)
			req := httptest.NewRequest(
				http.MethodPost,
				fmt.Sprintf("/events/%s/seats", tc.id),
				strings.NewReader(tc.body),
			)
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)
			assert.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestAddSeatsToEventHandlerSuccess(t *testing.T) {
	service := &MockSeatService{amount: 3}
	handler := NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /events/{eventID}/seats", handler.AddSeatsToEvent)

	req := httptest.NewRequest(
		http.MethodPost,
		"/events/5/seats",
		strings.NewReader(`{
		  "seats": [
			{
			  "number": "a1",
			  "price": 10
			},
			{
			  "number": "a2",
			  "price": 20
			},
			{
			  "number": "a3",
			  "price": 30
			}
		  ]
		}`,
		),
	)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, int64(3), service.amount)
	assert.Equal(t, http.StatusCreated, recorder.Code)
}

func TestGetSeatsByEventIDHandlerEmpty(t *testing.T) {
	service := &MockSeatService{seats: []Seat{}}
	handler := NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /events/{eventID}/seats", handler.GetSeatsByEventID)

	req := httptest.NewRequest(http.MethodGet, "/events/1/seats", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetSeatsByEventIDHandlerErrors(t *testing.T) {
	cases := []struct {
		name         string
		id           string
		expectedCode int
		error        error
	}{
		{
			name:         "invalid eventID",
			id:           "qwe",
			expectedCode: http.StatusBadRequest,
			error:        nil,
		},
		{
			name:         "validation error",
			id:           "0",
			expectedCode: http.StatusBadRequest,
			error:        ValidationError{Message: "eventID must be greater than zero"},
		},
		{
			name:         "internal server error",
			id:           "1",
			expectedCode: http.StatusInternalServerError,
			error:        errors.New("seats service error"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			service := &MockSeatService{err: tc.error}
			handler := NewHandler(service)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /events/{eventID}/seats", handler.GetSeatsByEventID)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/events/%s/seats", tc.id), nil)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)
			assert.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestGetSeatsByEventIDHandlerSuccess(t *testing.T) {
	service := &MockSeatService{
		seats: []Seat{
			{
				ID:        1,
				EventID:   1,
				Number:    "a1",
				Price:     100,
				CreatedAt: time.Time{},
			},
			{
				ID:        2,
				EventID:   1,
				Number:    "a2",
				Price:     150,
				CreatedAt: time.Time{},
			},
		},
	}

	handler := NewHandler(service)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /events/{eventID}/seats", handler.GetSeatsByEventID)
	req := httptest.NewRequest(http.MethodGet, "/events/1/seats", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}
