package sensors

import (
	"github.com/erickgreco/indoorgrid-system/pkg/response"
	"github.com/erickgreco/indoorgrid-system/pkg/validator"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) DHT11Handler(c *gin.Context) {
	var payload DHT11Payload

	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := validator.Validate.Struct(payload); err != nil {
		response.BadRequest(c, err)
		return
	}

	if err := h.service.checkDHT11(c.Request.Context(), payload); err != nil {
		response.InternalError(c)
		return
	}

	response.Created(c, response.SavedData)
}
