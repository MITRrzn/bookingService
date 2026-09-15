package booking

import (
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
