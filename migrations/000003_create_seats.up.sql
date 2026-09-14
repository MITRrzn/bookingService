CREATE TABLE seats
(
    id         BIGSERIAL PRIMARY KEY,
    event_id   BIGINT         NOT NULL,
    number     VARCHAR(50)    NOT NULL,
    price      NUMERIC(12, 2) NOT NULL,
    created_at TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_seats_event
        FOREIGN KEY (event_id)
            REFERENCES events (id)
            ON DELETE CASCADE,

    CONSTRAINT uq_seats_event_number
        UNIQUE (event_id, number),

    CONSTRAINT chk_seats_price
        CHECK (price >= 0)
);

CREATE INDEX idx_seats_event_id
    ON seats (event_id);