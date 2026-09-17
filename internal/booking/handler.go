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

	result, bookErr := h.service.ReserveSeat(r.Context(), eventId, seatId, req)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	encodeErr := json.NewEncoder(w).Encode(Response{
		Success:     "OK",
		BookingData: result,
	})

	if encodeErr != nil {
		log.Println("encode response error:", encodeErr)
	}
}

func (h *Handler) Confirm(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("bookingID")
	bookingId, parseErr := strconv.ParseInt(id, 10, 64)
	if parseErr != nil {
		log.Println("invalid bookingID id:", parseErr)
		helper.WriteErrorResponse(w, "invalid bookingID id", http.StatusBadRequest)
		return
	}

	var req ConfirmBookingRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Println("decode error:", err)
		helper.WriteErrorResponse(w, "incorrect json format", http.StatusBadRequest)
		return
	}

	result, confirmErr := h.service.ConfirmBooking(r.Context(), bookingId, req)
	if confirmErr != nil {
		var validationErr ValidationError
		var conflictErr ConflictError

		if errors.As(confirmErr, &validationErr) {
			helper.WriteErrorResponse(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		if errors.As(confirmErr, &conflictErr) {
			helper.WriteErrorResponse(w, conflictErr.Error(), http.StatusConflict)
			return
		}

		log.Println(confirmErr)
		helper.WriteErrorResponse(w, "failed confirm booking", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encodeErr := json.NewEncoder(w).Encode(Response{
		Success:     "OK",
		BookingData: result,
	})

	if encodeErr != nil {
		log.Println("encode response error:", encodeErr)
	}
}

func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("bookingID")
	bookingID, parseErr := strconv.ParseInt(id, 10, 64)
	if parseErr != nil {
		log.Println("invalid booking id:", parseErr)
		helper.WriteErrorResponse(w, "invalid booking id", http.StatusBadRequest)
		return
	}

	var req CancelBookingRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Println("decode error:", err)
		helper.WriteErrorResponse(w, "incorrect json format", http.StatusBadRequest)
		return
	}

	cancelErr := h.service.CancelBooking(r.Context(), bookingID, req)
	if cancelErr != nil {
		var conflictErr ConflictError
		var internalErr InternalError
		var validationErr ValidationError

		if errors.As(cancelErr, &validationErr) {
			helper.WriteErrorResponse(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		if errors.As(cancelErr, &conflictErr) {
			helper.WriteErrorResponse(w, conflictErr.Error(), http.StatusConflict)
			return
		}
		if errors.As(cancelErr, &internalErr) {
			helper.WriteErrorResponse(w, internalErr.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	var req ListBookingRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Println("decode error:", err)
		helper.WriteErrorResponse(w, "incorrect json format", http.StatusBadRequest)
		return
	}

	list, listErr := h.service.ListBookings(r.Context(), req.UserID)
	if listErr != nil {
		var validationErr ValidationError

		if errors.As(listErr, &validationErr) {
			helper.WriteErrorResponse(w, validationErr.Error(), http.StatusBadRequest)
			return
		}

		helper.WriteErrorResponse(w, "something goes wrong", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encodeErr := json.NewEncoder(w).Encode(ListResponse{
		Success:     "OK",
		BookingData: list,
	})

	if encodeErr != nil {
		log.Println("encode response error:", encodeErr)
	}
}
