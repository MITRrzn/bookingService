package booking

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewBookingRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) ReserveSeat(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
