-- +goose Up
CREATE TABLE bme280_readings (
    id UUID PRIMARY KEY,
    temperature NUMERIC(5,2) NOT NULL,
    humidity NUMERIC(5,2) NOT NULL,
    pressure NUMERIC(7,2) NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE bh1750_readings (
    id UUID PRIMARY KEY,
    lux NUMERIC(8,2) NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE soil_readings (
    id UUID PRIMARY KEY,
    moisture_percent NUMERIC(5,2) NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE soil_readings;
DROP TABLE bh1750_readings;
DROP TABLE bme280_readings;
