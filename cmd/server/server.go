package server

import (
	"github.com/erickgreco/indoorgrid-system/cmd/config"
	"github.com/erickgreco/indoorgrid-system/internal/camera/gopro"
	"github.com/erickgreco/indoorgrid-system/internal/mqtt"
	"github.com/erickgreco/indoorgrid-system/internal/sensors"
	"github.com/erickgreco/indoorgrid-system/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

const version string = "0.0.1"

type Server struct {
	router *gin.Engine
	db     *pgxpool.Pool
	cfg    config.Config
	camera *gopro.GoPro
	mqtt   *mqtt.Client
}

func New(db *pgxpool.Pool, cfg config.Config, camera *gopro.GoPro, mqttClient *mqtt.Client) *Server {
	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)

	s := &Server{
		router: r,
		db:     db,
		cfg:    cfg,
		camera: camera,
		mqtt:   mqttClient,
	}
	s.wire()
	return s
}

func (s *Server) registerRoutes() {
	api := s.router.Group("/v1")
	{
		api.GET("/health", s.healthCheckHandler)
	}

}

func (s *Server) Run() error {
	return s.router.Run(s.cfg.Port)
}

func (s *Server) wire() {
	sensorsRepo := sensors.NewRepo(s.db)
	sensorsService := sensors.NewService(sensorsRepo)

	if err := s.mqtt.Subscribe(sensorsService); err != nil {
		logger.Warn(logger.SubscribeErr, err)
		return
	}

	s.registerRoutes()
}
