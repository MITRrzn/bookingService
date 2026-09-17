package main

import (
	"bookingService/internal/database/psql"
	"bookingService/internal/repository/expired"
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
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
	bookingRepo := expired.NewExpiredRepo(db)

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			result, repoErr := bookingRepo.UpdateExpiredBookings(ctx)
			if repoErr != nil {
				log.Println(repoErr)
			} else {
				log.Printf("Update expired bookings: %v", result)
			}
		case <-ctx.Done():
			return
		}
	}
}
