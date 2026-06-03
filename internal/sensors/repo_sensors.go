package sensors

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewSensorsRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) insertDHT11(ctx context.Context, s DHT11) error {
	query := `
		INSERT INTO dht11_readings (id, temperature, humidity, unusual)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		s.ID,
		s.Temperature,
		s.Humidity,
		s.Unusual,
	)
	return err
}
