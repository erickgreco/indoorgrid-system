-- +goose Up
CREATE TABLE grow_cycles (
    id UUID PRIMARY KEY,
    plant_name TEXT NOT NULL,
    phase TEXT NOT NULL CHECK (phase IN ('germination', 'vegetative', 'flowering', 'harvest')),
    light_hours SMALLINT NOT NULL CHECK (light_hours BETWEEN 1 AND 24),
    dark_hours SMALLINT NOT NULL CHECK (dark_hours BETWEEN 0 AND 23),
    notes TEXT,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ,
    CONSTRAINT light_dark_hours_sum CHECK (light_hours + dark_hours = 24)
);

CREATE UNIQUE INDEX one_active_cycle
    ON grow_cycles ((ended_at IS NULL))
    WHERE ended_at IS NULL;

CREATE INDEX idx_grow_cycles_started_at
    ON grow_cycles (started_at DESC);

-- +goose Down
DROP TABLE IF EXISTS grow_cycles;
