-- +goose Up
CREATE TABLE soil_readings (
    id UUID PRIMARY KEY,
    grow_cycle_id UUID REFERENCES grow_cycles(id),
    moisture_percent NUMERIC(5,2) NOT NULL,
    sensor_at TIMESTAMPTZ NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_soil_readings_grow_cycle_id
    ON soil_readings (grow_cycle_id);

CREATE INDEX idx_soil_readings_sensor_at
    ON soil_readings (sensor_at DESC);

-- +goose Down
DROP TABLE IF EXISTS soil_readings;
