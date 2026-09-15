package seats

import "net/http"

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

func (h *Handler) GetSeatsByEventId(w http.ResponseWriter, r *http.Request) {

}
