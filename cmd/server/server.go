package server

import (
	"github.com/erickgreco/indoorgrid-system/cmd/config"
	"github.com/erickgreco/indoorgrid-system/internal/sensors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const version string = "0.0.1"

type Server struct {
	router *gin.Engine
	db     *pgxpool.Pool
	cfg    config.Config
}

func New(db *pgxpool.Pool, cfg config.Config) *Server {
	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)

	s := &Server{
		router: r,
		db:     db,
		cfg:    cfg,
	}
	s.wire()
	return s
}

func (s *Server) registerRoutes(sensors *sensors.Handler) {
	api := s.router.Group("/v1")
	{
		api.GET("/health", s.healthCheckHandler)

		sensorsGroup := api.Group("/sensors")
		{
			sensorsGroup.POST("/dht11", sensors.DHT11Handler)
		}
	}

}

func (s *Server) Run() error {
	return s.router.Run(s.cfg.Port)
}

func (s *Server) wire() {
	sensorsRepo := sensors.NewSensorsRepo(s.db)
	sensorsService := sensors.NewService(sensorsRepo)
	sensorsHandler := sensors.NewHandler(sensorsService)

	s.registerRoutes(sensorsHandler)

}
