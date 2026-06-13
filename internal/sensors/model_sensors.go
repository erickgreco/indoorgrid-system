package sensors

import (
	"time"

	"github.com/google/uuid"
)

const (
	minTemperature       = 0.00
	maxTemperature       = 45.00
	minHumidity          = 10.00
	maxHumidity          = 100.00
	minPressureHpa       = 300.00
	maxPressureHpa       = 1100.00
	minGasResistanceOhms = 50.00
	maxGasResistanceOhms = 500000.00
	minIlluminanceLux    = 0.11
	maxIlluminanceLux    = 100000.00
	minSoilPercent       = 0.00
	maxSoilPercent       = 100.00
)

type BME680Reading struct {
	ID                 uuid.UUID
	GrowCycleID        *uuid.UUID
	TemperatureCelcius float64
	HumidityPercent    float64
	PressureHpa        float64
	GasResistanceOhms  float64
	SensorAt           time.Time
	RecordedAt         time.Time
}

type BH1750Reading struct {
	ID             uuid.UUID
	GrowCycleID    *uuid.UUID
	IlluminanceLux float64
	SensorAt       time.Time
	RecordedAt     time.Time
}

type SoilReading struct {
	ID              uuid.UUID
	GrowCycleID     *uuid.UUID
	MoisturePercent float64
	SensorAt        time.Time
	RecordedAt      time.Time
}
