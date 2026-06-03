-- +goose Up
ALTER TABLE dht11_readings ADD COLUMN unusual BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE dht11_readings DROP COLUMN unusual;
