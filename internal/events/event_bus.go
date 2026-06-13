package events

import "github.com/erickgreco/indoorgrid-system/internal/sensors"

type EventBus struct {
	BME680 chan sensors.BME680Reading
	BH1750 chan sensors.BH1750Reading
	Soil   chan sensors.SoilReading
}

func New() *EventBus {
	return &EventBus{
		BME680: make(chan sensors.BME680Reading),
		BH1750: make(chan sensors.BH1750Reading),
		Soil:   make(chan sensors.SoilReading),
	}
}
