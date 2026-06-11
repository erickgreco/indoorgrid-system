-- +goose Up
CREATE TABLE grow_cycles (
    id UUID PRIMARY KEY,
    phase TEXT NOT NULL CHECK (phase IN ('germination', 'vegetative', 'flowering')),
    plant_name TEXT NOT NULL,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX one_active_cycle
    ON grow_cycles ((ended_at IS NULL))
    WHERE ended_at IS NULL;

CREATE INDEX idx_grow_cycles_started_at
    ON grow_cycles (started_at DESC);

-- +goose Down
DROP TABLE IF EXISTS grow_cycles;
