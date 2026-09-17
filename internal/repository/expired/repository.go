package expired

import (
	"context"
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewExpiredRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) UpdateExpiredBookings(ctx context.Context) (int64, error) {
	result, execErr := r.db.ExecContext(ctx, `UPDATE bookings SET status = 'expired' WHERE status = 'reserved' AND expires_at <= now()`)
	if execErr != nil {
		return 0, execErr
	}

	amount, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return amount, nil
}
