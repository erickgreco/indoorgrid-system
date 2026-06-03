package sensors

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) checkDHT11(ctx context.Context, sensor DHT11Payload) error {
	unusual := sensor.Temperature < dht11MinTemp || sensor.Temperature > dht11MaxTemp ||
		sensor.Humidity < dht11MinHum || sensor.Humidity > dht11MaxHum

	data := DHT11{
		ID:          uuid.New(),
		Temperature: sensor.Temperature,
		Humidity:    sensor.Humidity,
		Unusual:     unusual,
	}

	return s.repo.insertDHT11(ctx, data)
}
