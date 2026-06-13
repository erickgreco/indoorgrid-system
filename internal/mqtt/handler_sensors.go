package mqtt

import (
	"context"
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/erickgreco/indoorgrid-system/internal/sensors"
	"github.com/erickgreco/indoorgrid-system/pkg/logger"
)

func (c *Client) BME680Handler(service *sensors.Service) mqtt.MessageHandler {
	return func(_ mqtt.Client, m mqtt.Message) {
		var reading sensors.BME680Reading

		if err := json.Unmarshal(m.Payload(), &reading); err != nil {
			logger.Warn(logger.ParsePayloadErr, err)
			return
		}

		ctx := context.Background()

		if err := service.SaveBME680(ctx, reading); err != nil {
			logger.Warn(logger.SaveBME680Err, err)
			return
		}
	}
}

func (c *Client) BH1750Handler(service *sensors.Service) mqtt.MessageHandler {
	return func(_ mqtt.Client, m mqtt.Message) {
		var reading sensors.BH1750Reading

		if err := json.Unmarshal(m.Payload(), &reading); err != nil {
			logger.Warn(logger.ParsePayloadErr, err)
			return
		}

		ctx := context.Background()

		if err := service.SaveBH1750(ctx, reading); err != nil {
			logger.Warn(logger.SaveBH1750Err, err)
			return
		}
	}
}

func (c *Client) SoilHandler(service *sensors.Service) mqtt.MessageHandler {
	return func(_ mqtt.Client, m mqtt.Message) {
		var reading sensors.SoilReading

		if err := json.Unmarshal(m.Payload(), &reading); err != nil {
			logger.Warn(logger.ParsePayloadErr, err)
			return
		}

		ctx := context.Background()

		if err := service.SaveSoil(ctx, reading); err != nil {
			logger.Warn(logger.SaveSoilErr, err)
			return
		}
	}
}
