package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"spacebox-writer/internal/configs"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"spacebox-writer/internal/app"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

const (
	DefaultEnvFile = ".env"
	EnvFile        = "ENV_FILE"
)

func main() {
	//flag.StringVar(&configPath, "config", "config/local.yaml", "Config file path")
	//flag.Parse()

	// try to get .env file from Environments
	fileName, ok := os.LookupEnv(EnvFile)
	if !ok {
		fileName = DefaultEnvFile
	}

	// load environment variables based on .env file
	if err := godotenv.Load(fileName); err != nil {
		panic(err)
	}

	var cfg configs.Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg)

	application := app.New(cfg)

	startCtx, startCancel := context.WithTimeout(context.Background(), cfg.StartTimeout)
	defer startCancel()
	if err := application.Start(startCtx); err != nil {
		log.Fatal().Err(err).Msg("cannot start application") // nolint
	}

	log.Info().Msg("application started")

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quitCh

	stopCtx, stopCancel := context.WithTimeout(context.Background(), cfg.StopTimeout)
	defer stopCancel()

	if err := application.Stop(stopCtx); err != nil {
		log.Error().Err(err).Msg("cannot stop application")
	}

	log.Info().Msg("service is down")
}
