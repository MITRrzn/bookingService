package event

import "net/http"

func AddNewEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func GetActiveEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func GetEventById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
