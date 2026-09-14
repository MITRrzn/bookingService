package event

import (
	"bookingService/internal/helper"
	"bookingService/internal/structs"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Handler struct {
	service Service
}

func (h *Handler) CreateEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var event structs.CreateEventInput

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&event); err != nil {
			log.Println("decode error:", err)
			helper.WriteErrorResponse(w, "incorrect json format", http.StatusBadRequest)
			return
		}

		createdEvent, createErr := h.service.CreateEvent(r.Context(), event)

		if createErr != nil {
			if errors.Is(createErr, ValidationError) {
				helper.WriteErrorResponse(
					w,
					createErr.Error(),
					http.StatusBadRequest,
				)
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
}

// func (h *Handler) GetEvents() http.HandlerFunc {
// }

// func (h *Handler) GetEventByID() http.HandlerFunc {
// }
