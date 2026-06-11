-- +goose Up
CREATE TABLE indoor_sensors (
    id UUID PRIMARY KEY,
    temperature_celsius NUMERIC(5,2),
    humidity_percent NUMERIC(5, 2),
    pressure_hpa NUMERIC(7,2),
    gas_resistance_ohms NUMERIC(12, 2),
    illuminance_lux NUMERIC(8,2),
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_indoor_sensors_recorded_at 
    ON indoor_sensors (recorded_at DESC);

-- +goose Down
DROP TABLE IF EXISTS indoor_sensors;