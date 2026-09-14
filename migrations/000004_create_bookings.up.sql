CREATE TABLE bookings
(
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT      NOT NULL,
    seat_id      BIGINT      NOT NULL,
    status       VARCHAR(20) NOT NULL,
    reserved_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at   TIMESTAMPTZ NOT NULL,
    confirmed_at TIMESTAMPTZ NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_bookings_user
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_bookings_seat
        FOREIGN KEY (seat_id)
            REFERENCES seats (id)
            ON DELETE CASCADE,

    CONSTRAINT chk_bookings_status
        CHECK (
            status IN (
                       'reserved',
                       'confirmed',
                       'cancelled',
                       'expired'
                )
            ),

    CONSTRAINT chk_bookings_expiration
        CHECK (expires_at > reserved_at)
);

CREATE INDEX idx_bookings_user_id
    ON bookings (user_id);

CREATE INDEX idx_bookings_seat_id
    ON bookings (seat_id);

CREATE INDEX idx_bookings_status
    ON bookings (status);

CREATE INDEX idx_bookings_expiration
    ON bookings (status, expires_at);

CREATE UNIQUE INDEX uq_bookings_confirmed_seat
    ON bookings (seat_id)
    WHERE status = 'confirmed';