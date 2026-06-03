-- +goose Up
CREATE TABLE dht11_readings (
    id UUID PRIMARY KEY,
    temperature NUMERIC(5,2) NOT NULL,
    humidity NUMERIC(5,2) NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE dht11_readings;
