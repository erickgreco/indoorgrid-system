-- +goose Up
CREATE TABLE bme680_readings (
    id UUID PRIMARY KEY,
    grow_cycle_id UUID REFERENCES grow_cycles(id),
    temperature_celcius NUMERIC(5,2) NOT NULL,
    humidity_percent NUMERIC(5,2) NOT NULL,
    pressure_hpa NUMERIC(7,2) NOT NULL,
    gas_resistance_ohms NUMERIC(12,2) NOT NULL,
    sensor_at TIMESTAMPTZ NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_bme680_readings_grow_cycle_id
    ON bme680_readings (grow_cycle_id);

CREATE INDEX idx_bme680_readings_sensor_at
    ON bme680_readings (sensor_at DESC);

-- +goose Down
DROP TABLE IF EXISTS bme680_readings;
