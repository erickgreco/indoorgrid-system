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
	ID                 uuid.UUID  `json:"id"`
	GrowCycleID        *uuid.UUID `json:"grow_cycle_id"`
	TemperatureCelcius float64    `json:"temperature_celcius"`
	HumidityPercent    float64    `json:"humidity_percent"`
	PressureHpa        float64    `json:"pressure_hpa"`
	GasResistanceOhms  float64    `json:"gas_resistance_ohms"`
	SensorAt           time.Time  `json:"sensor_at"`
	RecordedAt         time.Time  `json:"recorded_at"`
}

type BH1750Reading struct {
	ID             uuid.UUID  `json:"id"`
	GrowCycleID    *uuid.UUID `json:"grow_cycle_id"`
	IlluminanceLux float64    `json:"illuminance_lux"`
	SensorAt       time.Time  `json:"sensor_at"`
	RecordedAt     time.Time  `json:"recorded_at"`
}

type SoilReading struct {
	ID              uuid.UUID  `json:"id"`
	GrowCycleID     *uuid.UUID `json:"grow_cycle_id"`
	MoisturePercent float64    `json:"moisture_percent"`
	SensorAt        time.Time  `json:"sensor_at"`
	RecordedAt      time.Time  `json:"recorded_at"`
}
