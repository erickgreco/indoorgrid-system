-- +goose Up
CREATE TABLE bh1750_readings (
    id UUID PRIMARY KEY,
    grow_cycle_id UUID REFERENCES grow_cycles(id),
    illuminance_lux NUMERIC(8,2) NOT NULL,
    sensor_at TIMESTAMPTZ NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bh1750_readings_grow_cycle_id
    ON bh1750_readings (grow_cycle_id);

CREATE INDEX idx_bh1750_readings_sensor_at
    ON bh1750_readings (sensor_at DESC);

-- +goose Down
DROP TABLE IF EXISTS bh1750_readings;
