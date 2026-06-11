-- +goose Up
ALTER TABLE indoor_sensors
    ADD COLUMN grow_cycle_id UUID REFERENCES grow_cycles(id);

CREATE INDEX idx_indoor_sensors_grow_cycle_id
    ON indoor_sensors (grow_cycle_id);

-- +goose Down
ALTER TABLE indoor_sensors
    DROP COLUMN IF EXISTS grow_cycle_id;
