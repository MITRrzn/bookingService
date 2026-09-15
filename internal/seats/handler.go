package seats

import (
	"bookingService/internal/helper"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *SeatService
}

func NewHandler(service *SeatService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) AddSeatsToEvent(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetSeatsByEventID(w http.ResponseWriter, r *http.Request) {
	eventID := r.PathValue("eventID")
	id, parseErr := strconv.ParseInt(eventID, 10, 64)
	if parseErr != nil {
		log.Println("invalid event id:", parseErr)
		helper.WriteErrorResponse(w, "invalid event id", http.StatusBadRequest)
		return
	}

	result, err := h.service.GetSeatsByEventID(r.Context(), id)
	if err != nil {
		helper.WriteErrorResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encodeErr := json.NewEncoder(w).Encode(SeatsResponse{
		Success: "OK",
		Seats:   result,
	})

	if encodeErr != nil {
		log.Println("encode response error:", encodeErr)
	}
}
