package main

import (
	"log"
	"time"

	"github.com/erickgreco/indoorgrid-system/cmd/config"
	"github.com/erickgreco/indoorgrid-system/cmd/server"

	"github.com/erickgreco/indoorgrid-system/internal/camera/gopro"
	"github.com/erickgreco/indoorgrid-system/internal/db"
	"github.com/erickgreco/indoorgrid-system/pkg/env"
)

func main() {
	go func() {
		_, err := gopro.BluetoothConn()
		if err != nil {
			log.Println(err)
		}
	}()

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
	log.Println("database connection established")

	srv := server.New(dbpool, cfg)
	if err := srv.Run(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
