package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/bro-n-bro/spacebox-writer/internal/app"
)

const (
	DefaultEnvFile = ".env.local"
	EnvFile        = "ENV_FILE"

	msgCannotStartApplication = "cannot start application"
	msgServiceIsDown          = "service is down"
	msgCannotStopApplication  = "cannot stop application"
	msgApplicationStarted     = "application started"
)

// init is a function for setting up logger before main function
func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	fileName, ok := os.LookupEnv(EnvFile)
	if !ok {
		fileName = DefaultEnvFile
	}

	if err := godotenv.Load(fileName); err != nil {
		panic(err)
	}

	var cfg app.Config
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	application := app.New(cfg, l)

	startCtx, startCancel := context.WithTimeout(context.Background(), cfg.StartTimeout)
	defer startCancel()
	if err := application.Start(startCtx); err != nil {
		log.Fatal().Err(err).Msg(msgCannotStartApplication) // nolint
	}

	log.Info().Msg(msgApplicationStarted)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quitCh

	stopCtx, stopCancel := context.WithTimeout(context.Background(), cfg.StopTimeout)
	defer stopCancel()

	if err := application.Stop(stopCtx); err != nil {
		log.Error().Err(err).Msg(msgCannotStopApplication)
	}

	log.Info().Msg(msgServiceIsDown)
	// os.Exit()
}
