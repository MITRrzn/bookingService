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
	eventID := r.PathValue("eventID")
	id, parseErr := strconv.ParseInt(eventID, 10, 64)
	if parseErr != nil {
		log.Println("invalid event id:", parseErr)
		helper.WriteErrorResponse(w, "invalid event id", http.StatusBadRequest)
		return
	}

	var request AddSeatsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		helper.WriteErrorResponse(w, "incorrect input data", http.StatusBadRequest)
		return
	}

	amount, addErr := h.service.AddSeatsToEvent(r.Context(), request.Seats, id)
	if addErr != nil {
		helper.WriteErrorResponse(w, "failed to add seats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encodeErr := json.NewEncoder(w).Encode(PostSeatsResponse{
		Success: "OK",
		Amount:  amount,
	})
	if encodeErr != nil {
		log.Println("encode response error:", encodeErr)
	}
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

	encodeErr := json.NewEncoder(w).Encode(GetSeatsResponse{
		Success: "OK",
		Seats:   result,
	})

	if encodeErr != nil {
		log.Println("encode response error:", encodeErr)
	}
}
