package server

import (
	"github.com/erickgreco/indoorgrid-system/cmd/config"
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
	s.registerRoutes()
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
