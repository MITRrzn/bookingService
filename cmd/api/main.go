package main

import (
	"bookingService/internal/database/psql"
	"bookingService/internal/event"
	eventRepository "bookingService/internal/repository/event"
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

	eventRepo := eventRepository.NewRepository(db)
	eventService := event.NewService(eventRepo)
	eventHandler := event.NewHandler(eventService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /events", eventHandler.CreateEvent())
	//mux.HandleFunc("GET /events", event.GetActiveEvents())
	mux.HandleFunc("GET /events/{id}", eventHandler.GetEventByID())

	port := os.Getenv("APP_PORT")
	log.Println("Starting server at port", port)
	server := &http.Server{
		Addr:              fmt.Sprint(":", port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
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
