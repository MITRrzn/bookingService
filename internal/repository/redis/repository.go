package redis

import (
	"bookingService/internal/booking"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type ReservationStore struct {
	client *redis.Client
}

func NewReservationStore(client *redis.Client) *ReservationStore {
	return &ReservationStore{
		client: client,
	}
}

func (r *ReservationStore) Reserve(ctx context.Context, input booking.Input, ttl time.Duration) (bool, error) {
	key := fmt.Sprintf("reservation:%d:%d", input.EventID, input.SeatID)

	reserved, err := r.client.SetNX(ctx, key, input.UserID, ttl).Result()

	if err != nil {
		return false, err
	}

	return reserved, nil
}

func (r *ReservationStore) DeleteReserve(ctx context.Context, seatID int64) error {
	key := fmt.Sprintf("reservation:%d", seatID)

	_, err := r.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}
