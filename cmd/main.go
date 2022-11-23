package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"spacebox-writer/internal/app"

	"github.com/jinzhu/configor"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var configPath string

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	flag.StringVar(&configPath, "config", "config/local.yaml", "Config file path")
	flag.Parse()

	var cfg app.Config
	if err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&cfg, configPath); err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

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
