package events

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	bus *EventBus
}

func NewHandler(bus *EventBus) *Handler {
	return &Handler{
		bus: bus,
	}
}

func (h *Handler) LiveHandler(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		select {
		case reading := <-h.bus.BME680:
			data, _ := json.Marshal(reading)
			c.SSEvent("bme680", string(data))
			return true
		case reading := <-h.bus.BH1750:
			data, _ := json.Marshal(reading)
			c.SSEvent("bh1750", string(data))
			return true
		case reading := <-h.bus.Soil:
			data, _ := json.Marshal(reading)
			c.SSEvent("soil", string(data))
			return true
		case <-c.Request.Context().Done():
			return false
		}
	})
}
