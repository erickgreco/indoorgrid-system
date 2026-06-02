package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": version,
	})
}
