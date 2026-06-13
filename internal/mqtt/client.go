package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/erickgreco/indoorgrid-system/cmd/config"
	"github.com/erickgreco/indoorgrid-system/internal/events"
	"github.com/erickgreco/indoorgrid-system/internal/sensors"
	"github.com/erickgreco/indoorgrid-system/pkg/logger"
)

const (
	bme680tag = "indoorgrid/sensors/bme680"
	bh1750tag = "indoorgrid/sensors/bh1750"
	soiltag   = "indoorgrid/sensors/soil"
)

type Client struct {
	conn mqtt.Client
}

func New(cfg config.Config) (*Client, error) {
	opts := mqtt.NewClientOptions()

	opts.AddBroker(cfg.MQTTBroker)
	opts.SetClientID(cfg.MQTTClientID)
	opts.SetAutoReconnect(true)
	opts.SetCleanSession(false)

	conn := mqtt.NewClient(opts)
	if token := conn.Connect(); token.Wait() && token.Error() != nil {
		return nil, logger.Error(logger.MQTTConnErr, token.Error())
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Subscribe(service *sensors.Service, bus *events.EventBus) error {
	topics := map[string]mqtt.MessageHandler{
		bme680tag: c.BME680Handler(service, bus),
		bh1750tag: c.BH1750Handler(service, bus),
		soiltag:   c.SoilHandler(service, bus),
	}

	for topic, handler := range topics {
		token := c.conn.Subscribe(topic, 1, handler)
		token.Wait()

		if err := token.Error(); err != nil {
			return logger.Error(logger.SubscribeErr, err, logger.Topic, topic)
		}
	}

	return nil
}
