package main

import (
	"log"
	"time"

	"github.com/erickgreco/indoorgrid-system/cmd/config"
	"github.com/erickgreco/indoorgrid-system/cmd/server"

	"github.com/erickgreco/indoorgrid-system/internal/camera/gopro"
	"github.com/erickgreco/indoorgrid-system/internal/db"
	"github.com/erickgreco/indoorgrid-system/pkg/env"
	"github.com/erickgreco/indoorgrid-system/pkg/logger"
)

func main() {
	logger.Init(env.GetString("LOG_LEVEL", "debug"))

	cfg := config.Load()

	dbcfg := db.Config{
		MaxConns:              10,
		MinConns:              5,
		MaxConnLifeTime:       time.Hour,
		MaxConnIdleTime:       10 * time.Minute,
		HealthCheckPeriod:     time.Minute,
		MaxConnLifeTimeJitter: 5 * time.Minute,
	}

	dbpool, err := db.Connect(env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5433/indoorgrid-system?sslmode=disable"), dbcfg)
	if err != nil {
		log.Fatalf("db connection err: %v", err)
	}
	defer dbpool.Close()
	logger.Info("database connection established")

	goPro := gopro.New()

	go func() {
		if err := goPro.BleConn(); err != nil {
			logger.Error(logger.DeviceConnErr, err)
		}
	}()

	srv := server.New(dbpool, cfg, goPro)
	if err := srv.Run(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
