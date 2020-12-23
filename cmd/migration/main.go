package main

import (
	"context"
	"flag"
	"go-api-template/internal/config"
	"go-api-template/pkg/db"
	"go-api-template/pkg/log"
	"os"
	"time"

	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
)

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")
var flagMigrations = flag.String("migrations", "./cmd/migrations", "path to the migrations folder")

func main() {
	flag.Parse()

	logger := log.Setup()

	cfg, err := config.New(*flagConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	db, err := db.Connect(ctx, logger.With().Str("layer", "database").Logger(), cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	err = migrations.Run(db, *flagMigrations, os.Args)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
}
