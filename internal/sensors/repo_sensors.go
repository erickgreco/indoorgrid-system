package sensors

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) insertBME680(ctx context.Context, reading BME680Reading) error {
	query := `
		INSERT INTO bme680_readings (id, grow_cycle_id, temperature_celcius, humidity_percent, pressure_hpa,
		gas_resistance_ohms, sensor_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		reading.ID,
		reading.GrowCycleID,
		reading.TemperatureCelcius,
		reading.HumidityPercent,
		reading.PressureHpa,
		reading.GasResistanceOhms,
		reading.SensorAt,
	)

	return err
}

func (r *Repo) insertBH1750(ctx context.Context, reading BH1750Reading) error {
	query := `
		INSERT INTO bh1750_readings (id, grow_cycle_id, illuminance_lux, sensor_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		reading.ID,
		reading.GrowCycleID,
		reading.IlluminanceLux,
		reading.SensorAt,
	)

	return err
}

func (r *Repo) insertSoil(ctx context.Context, reading SoilReading) error {
	query := `
		INSERT INTO soil_readings (id, grow_cycle_id, moisture_percent, sensor_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		reading.ID,
		reading.GrowCycleID,
		reading.MoisturePercent,
		reading.SensorAt,
	)

	return err
}

func (r *Repo) getActiveCycleID(ctx context.Context) (*uuid.UUID, error) {
	var id uuid.UUID

	query := `
		SELECT id 
		FROM grow_cycles
		WHERE ended_at IS NULL LIMIT 1
	`

	err := r.db.QueryRow(
		ctx,
		query,
	).Scan(
		&id,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &id, nil
}
