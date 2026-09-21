package event

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockEventService struct {
	event     Event
	eventList []Event
	error     error
}

func (m *mockEventService) CreateEvent(ctx context.Context, input CreateEventInput) (Event, error) {
	return m.event, m.error
}

func (m *mockEventService) GetActiveEvents(ctx context.Context) ([]Event, error) {
	return m.eventList, m.error
}

func (m *mockEventService) GetEventByID(ctx context.Context, id int64) (Event, error) {
	return m.event, m.error
}

func TestCreateEventHandlerErrors(t *testing.T) {
	cases := []struct {
		name         string
		body         string
		error        error
		expectedCode int
	}{
		{
			name:         "invalid json",
			error:        nil,
			body:         `"name": "test event", "starts_at" : "2026-05-12 11:12:13"}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "validation error",
			body:         `"name": "ab", "starts_at" : "2026-05-12 11:12:13"}`,
			error:        nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "internal error",
			body:         `{"name": "test", "starts_at" : "2026-05-12 11:12:13"}`,
			error:        errors.New("internal error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	service := &mockEventService{}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.error != nil {
				service = &mockEventService{
					error: tc.error,
				}
			}
			req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(tc.body))

			handler := NewHandler(service)
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			handler.CreateEvent(recorder, req)
			assert.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestCreateEventHandlerSuccess(t *testing.T) {
	req := httptest.NewRequest(
		http.MethodPost,
		"/events",
		strings.NewReader(`{"name": "test event", "starts_at" : "2026-05-12 11:12:13"}`))

	service := &mockEventService{
		event: Event{ID: 1, Name: "test event"},
	}
	handler := NewHandler(service)

	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.GetEvents(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetEventByIDHandlerErrors(t *testing.T) {
	cases := []struct {
		name         string
		id           any
		error        error
		expectedCode int
	}{
		{
			name:         "event not found",
			id:           1,
			error:        sql.ErrNoRows,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "invalid event id",
			id:           "qwerty",
			error:        nil,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "internal error",
			id:           1,
			error:        errors.New("error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	service := &mockEventService{}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.error != nil {
				service = &mockEventService{
					error: tc.error,
				}
			}

			handler := NewHandler(service)

			mux := http.NewServeMux()
			mux.HandleFunc("GET /events/{id}", handler.GetEventByID)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/events/%v", tc.id), nil)
			recorder := httptest.NewRecorder()

			mux.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code)
		})
	}
}

func TestGetByIDHandlerSuccess(t *testing.T) {
	service := &mockEventService{
		event: Event{ID: 1, Name: "test event"},
	}

	handler := NewHandler(service)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /events/{id}", handler.GetEventByID)

	req := httptest.NewRequest(http.MethodGet, "/events/1", nil)
	recorder := httptest.NewRecorder()

	mux.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetEventsHandlerSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/events", nil)

	service := &mockEventService{}
	handler := NewHandler(service)

	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.GetEvents(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetEventsHandlerNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/events", nil)

	service := &mockEventService{
		eventList: []Event{},
	}
	handler := NewHandler(service)

	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.GetEvents(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestGetEventsHandlerError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/events", nil)

	service := &mockEventService{
		error: errors.New("error"),
	}
	handler := NewHandler(service)

	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.GetEvents(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}
