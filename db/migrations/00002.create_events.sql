BEGIN;

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    player_id BIGINT REFERENCES players(id),
    game_id BIGINT,
    type TEXT NOT NULL,
    amount BIGINT,
    currency TEXT,
    has_won BOOLEAN,
    created_at TIMESTAMPTZ NOT NULL,
    amount_eur BIGINT,
    description TEXT
);

COMMIT;