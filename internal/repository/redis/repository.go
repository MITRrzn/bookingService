package redis

import (
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

func (r *ReservationStore) Reserve(ctx context.Context, eventID int64, seatID int64, userID int64, ttl time.Duration) (bool, error) {
	key := fmt.Sprintf("reservation:%d:%d", eventID, seatID)

	reserved, err := r.client.SetNX(ctx, key, userID, ttl).Result()

	if err != nil {
		return false, err
	}

	return reserved, nil
}
