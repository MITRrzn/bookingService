package main

import (
	"bookingService/internal/booking"
	"bookingService/internal/database/psql"
	"bookingService/internal/event"
	eventRepository "bookingService/internal/repository/event"
	seatsRepository "bookingService/internal/repository/seat"
	"bookingService/internal/seats"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := psql.PsqlConnect()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if dbCloseErr := db.Close(); dbCloseErr != nil {
			log.Printf("db close error: %v", dbCloseErr)
		}
	}()

	mux := http.NewServeMux()

	eventRepo := eventRepository.NewEventRepo(db)
	eventService := event.NewService(eventRepo)
	eventHandler := event.NewHandler(eventService)
	mux.HandleFunc("POST /events", eventHandler.CreateEvent)
	mux.HandleFunc("GET /events", eventHandler.GetEvents)
	mux.HandleFunc("GET /events/{id}", eventHandler.GetEventByID)

	seatsRepo := seatsRepository.NewSeatRepo(db)
	seatsService := seats.NewService(seatsRepo)
	seatsHandler := seats.NewHandler(seatsService)

	mux.HandleFunc("POST /events/{eventID}/seats", seatsHandler.AddSeatsToEvent)
	mux.HandleFunc("GET /events/{eventID}/seats", seatsHandler.GetSeatsByEventID)

	bookingRepo := booking.NewBookingRepo(db)
	bookingService := booking.NewService(bookingRepo)
	bookingHandler := booking.NewHandler(bookingService)
	mux.HandleFunc("POST /events/{eventID}/seats/{seatID}/reserve", bookingHandler.ReserveSeat)

	port := os.Getenv("APP_PORT")
	log.Println("Starting server at port", port)
	server := &http.Server{
		Addr:              fmt.Sprint(":", port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		listenErr := server.ListenAndServe()
		if listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			log.Fatal(listenErr)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down server")
	shutdownErr := server.Shutdown(shutdownCtx)
	if shutdownErr != nil {
		log.Printf("server shutdown error: %v", shutdownErr)
	}
	log.Println("Server stopped")
}
