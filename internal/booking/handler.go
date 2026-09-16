package booking

import (
	"bookingService/internal/helper"
	"database/sql"
	"encoding/json"
	"errors"
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

	seat := r.PathValue("seatID")
	seatId, parseErr := strconv.ParseInt(seat, 10, 64)
	if parseErr != nil {
		log.Println("invalid event id:", parseErr)
		helper.WriteErrorResponse(w, "invalid seat id", http.StatusBadRequest)
		return
	}

	var req ReserveSeatsRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Println("decode error:", err)
		helper.WriteErrorResponse(w, "incorrect json format", http.StatusBadRequest)
		return
	}

	bookErr := h.service.ReserveSeat(r.Context(), eventId, seatId, req)
	if errors.Is(bookErr, sql.ErrNoRows) {
		helper.WriteErrorResponse(w, "data not found", http.StatusNotFound)
		return
	}
	if bookErr != nil {
		var notFoundErr NotFoundError
		var validationErr ValidationError
		var internalErr InternalError
		var conflictErr ConflictError

		if errors.As(bookErr, &internalErr) {
			helper.WriteErrorResponse(w, internalErr.Error(), http.StatusInternalServerError)
			return
		}

		if errors.As(bookErr, &conflictErr) {
			helper.WriteErrorResponse(w, conflictErr.Error(), http.StatusConflict)
			return
		}

		if errors.As(bookErr, &notFoundErr) {
			helper.WriteErrorResponse(w, notFoundErr.Error(), http.StatusNotFound)
			return
		}
		if errors.As(bookErr, &validationErr) {
			helper.WriteErrorResponse(w, validationErr.Error(), http.StatusBadRequest)
			return
		}
	}
}
