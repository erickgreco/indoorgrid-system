package sensors

import "github.com/jackc/pgx/v5/pgxpool"

type Repo struct {
	db *pgxpool.Pool
}

func NewSensorsRepo(db *pgxpool.Pool) *Repo {
	return &Repo{
		db: db,
	}
}
