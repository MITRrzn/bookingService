package seat

import (
	"bookingService/internal/seats"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) GetSeatByID(ctx context.Context, seatID int64) (seats.Seat, error) {
	var model seats.Seat

	err := r.db.QueryRowContext(
		ctx,
		`
            SELECT id, event_id, number, price, created_at
            FROM seats
            WHERE id = $1
        `,
		seatID,
	).Scan(
		&model.ID,
		&model.EventID,
		&model.Number,
		&model.Price,
		&model.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return seats.Seat{}, sql.ErrNoRows
	}

	if err != nil {
		return seats.Seat{}, fmt.Errorf("get seat by id: %w", err)
	}

	return model, nil
}

func NewSeatRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) AddSeats(ctx context.Context, seats []seats.SeatInput, eventID int64) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}

	defer tx.Rollback()

	stmt, err := tx.PrepareContext(
		ctx,
		`
            INSERT INTO seats (event_id, number, price)
            VALUES ($1, $2, $3)
        `,
	)
	if err != nil {
		return 0, fmt.Errorf("prepare insert seats: %w", err)
	}
	defer stmt.Close()

	amount := 0
	for _, item := range seats {
		_, err = stmt.ExecContext(
			ctx,
			eventID,
			item.Number,
			item.Price,
		)
		if err != nil {
			return 0, fmt.Errorf("insert seat %s: %w", item.Number, err)
		}
		amount++
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}

	return int64(amount), nil
}

func (r *Repository) GetSeatsByEventID(ctx context.Context, eventID int64) ([]seats.Seat, error) {
	rows, queryErr := r.db.QueryContext(
		ctx,
		`
		SELECT id, event_id, number, price, created_at FROM seats
		WHERE event_id = $1
		`,
		eventID,
	)

	if queryErr != nil {
		return nil, queryErr
	}
	defer rows.Close()

	seatsList := make([]seats.Seat, 0)
	for rows.Next() {
		var model seats.Seat
		scanErr := rows.Scan(&model.ID, &model.EventID, &model.Number, &model.Price, &model.CreatedAt)
		if scanErr != nil {
			return nil, scanErr
		}

		seatsList = append(seatsList, model)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, rowsErr
	}

	return seatsList, nil
}
