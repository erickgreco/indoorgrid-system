package sensors

import (
	"time"

	"github.com/google/uuid"
)

type BME280 struct {
	ID          uuid.UUID
	Temperature float64
	Humidity    float64
	Pressure    float64
	RecordedAt  time.Time
}

type BH1750 struct {
	ID         uuid.UUID
	Lux        float64
	RecordedAt time.Time
}

type SoilSensor struct {
	ID              uuid.UUID
	MoisturePercent float64
	RecordedAt      time.Time
}
