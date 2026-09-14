package event

import (
	"bookingService/internal/helper"
	"bookingService/internal/structs"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var event structs.CreateEventInput

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&event); err != nil {
		log.Println("decode error:", err)
		helper.WriteErrorResponse(w, "incorrect json format", http.StatusBadRequest)
		return
	}

	createdEvent, createErr := h.service.CreateEvent(r.Context(), event)

	if createErr != nil {
		var validationErr structs.ValidationError
		if errors.As(createErr, &validationErr) {
			helper.WriteErrorResponse(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		log.Println("create event error:", createErr)
		helper.WriteErrorResponse(w, "incorrect json format", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encodeErr := json.NewEncoder(w).Encode(structs.SuccessResponse{
		Success: "created",
		Event:   createdEvent,
	})
	if encodeErr != nil {
		log.Println(encodeErr)
		return
	}
}

func (h *Handler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	request := r.PathValue("id")
	id, parseErr := strconv.ParseInt(request, 10, 64)
	if parseErr != nil {
		log.Println("invalid event id:", parseErr)
		helper.WriteErrorResponse(w, "invalid event id", http.StatusBadRequest)
		return
	}

	event, getEventErr := h.service.GetEventByID(r.Context(), id)
	if errors.Is(getEventErr, sql.ErrNoRows) {
		helper.WriteErrorResponse(w, getEventErr.Error(), http.StatusNotFound)
		return
	}
	if getEventErr != nil {
		helper.WriteErrorResponse(w, getEventErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encodeErr := json.NewEncoder(w).Encode(structs.SuccessResponse{
		Success: "OK",
		Event:   event,
	})
	if encodeErr != nil {
		log.Println(encodeErr)
		return
	}
}

func (h *Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	events, getEventsErr := h.service.GetEvents(r.Context())
	if getEventsErr != nil {
		helper.WriteErrorResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encodeErr := json.NewEncoder(w).Encode(structs.EventsResponse{
		Success: "OK",
		Events:  events,
	})

	if encodeErr != nil {
		log.Println("encode response error:", encodeErr)
	}
}
