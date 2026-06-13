-- +goose Up
CREATE TABLE watering_events (
    id UUID PRIMARY KEY,
    grow_cycle_id UUID NOT NULL REFERENCES grow_cycles(id),
    watered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    fertilizer TEXT,
    quantity_ml NUMERIC(8,2),
    water_ph NUMERIC(4,2) CHECK (water_ph BETWEEN 0 AND 14),
    notes TEXT
);

CREATE INDEX idx_watering_events_grow_cycle_id
    ON watering_events (grow_cycle_id);

CREATE INDEX idx_watering_events_watered_at 
    ON watering_events (watered_at DESC);

-- +goose Down
DROP TABLE IF EXISTS watering_events;
