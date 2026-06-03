package sensors

import (
	"time"

	"github.com/google/uuid"
)

const (
	dht11MinTemp = 10.0
	dht11MaxTemp = 40.0
	dht11MinHum  = 20.0
	dht11MaxHum  = 90.0
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

type DHT11Payload struct {
	Temperature float64 `json:"temperature" binding:"required"`
	Humidity    float64 `json:"humidity" binding:"required"`
}

type DHT11 struct {
	ID          uuid.UUID
	Temperature float64
	Humidity    float64
	Unusual     bool
	RecordedAt  time.Time
}
