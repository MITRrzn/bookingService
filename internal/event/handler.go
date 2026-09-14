package event

import "net/http"

type Handler struct {
	service Service
}

func (h *Handler) CreateEvent() http.HandlerFunc {
}
func (h *Handler) GetEvents() http.HandlerFunc {
}
func (h *Handler) GetEventByID() http.HandlerFunc {
}
