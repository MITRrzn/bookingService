package booking

import (
	"bookingService/internal/helper"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	service *BookingService
}

func NewHandler(service *BookingService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ReserveSeat(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("eventID")
	eventId, parseErr := strconv.ParseInt(id, 10, 64)
	if parseErr != nil {
		log.Println("invalid event id:", parseErr)
		helper.WriteErrorResponse(w, "invalid event id", http.StatusBadRequest)
		return
	}

	seatId := r.PathValue("seatID")

	var req ReserveSeatsRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Println("decode error:", err)
		helper.WriteErrorResponse(w, "incorrect json format", http.StatusBadRequest)
		return
	}
}
