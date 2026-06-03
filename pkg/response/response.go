package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SavedData = "data successfully saved"
)

type Respose struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func Ok(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Respose{Data: data})
}

func Created(c *gin.Context, message string) {
	c.JSON(http.StatusCreated, Respose{Message: message})
}

func BadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, Respose{Error: err.Error()})
}

func InternalError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Respose{Error: "internal server error"})
}
