package app

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/hexy-dev/spacebox-writer/adapter/broker"
	ch "github.com/hexy-dev/spacebox-writer/adapter/clickhouse"
	"github.com/hexy-dev/spacebox-writer/adapter/metrics"
	"github.com/hexy-dev/spacebox-writer/adapter/mongo"
	"github.com/hexy-dev/spacebox-writer/consts"
	"github.com/hexy-dev/spacebox-writer/internal/rep"
	"github.com/hexy-dev/spacebox-writer/models"
	"github.com/hexy-dev/spacebox-writer/modules"
)

type (
	App struct {
		log  *zerolog.Logger
		cmps []cmp
		cfg  Config
	}

	cmp struct {
		Service rep.Lifecycle
		Name    string
	}
)

func New(cfg Config, l zerolog.Logger) *App {
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = 0
	}

	l = l.Level(level).With().Str("cmp", "app").Logger()
	return &App{
		log:  &l,
		cfg:  cfg,
		cmps: []cmp{},
	}
}

func (a *App) Start(ctx context.Context) error {
	a.log.Info().Msg("starting app")
	a.log.Info().
		Uint8("log_level", uint8(a.log.GetLevel())).
		Str("log_level_text", a.cfg.LogLevel).
		Msg("logger")

	clickhouse := ch.New(a.cfg.Clickhouse, *a.log)
	m := mongo.New(a.cfg.Mongo, *a.log)
	brk := broker.New(a.cfg.Broker, clickhouse, m, *a.log)
	mods := modules.New(a.cfg.Modules, clickhouse, *a.log, brk)
	metrics := metrics.New(a.cfg.Metrics, clickhouse, *a.log)

	a.cmps = append(
		a.cmps,
		cmp{clickhouse, "clickhouse"},
		cmp{mods, "modules"},
		cmp{m, "mongo"},
		cmp{brk, "broker"},
		cmp{metrics, "metrics"},
	)

	okCh, errCh := make(chan struct{}), make(chan error)
	go func() {
		for _, c := range a.cmps {
			a.log.Info().Msgf("%v is starting", c.Name)
			if err := c.Service.Start(ctx); err != nil {
				a.log.Error().Err(err).Msgf(consts.FmtCannotStart, c.Name)
				errCh <- errors.Wrapf(err, consts.FmtCannotStart, c.Name)
			}
		}
		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return models.ErrStartTimeout
	case err := <-errCh:
		return err
	case <-okCh:
		return nil
	}
}

func (a *App) Stop(ctx context.Context) error {
	a.log.Info().Msg("shutting down service...")

	okCh, errCh := make(chan struct{}), make(chan error)
	go func() {
		for _, c := range a.cmps {
			a.log.Info().Msgf("stopping %q...", c.Name)
			if err := c.Service.Stop(ctx); err != nil {
				a.log.Error().Err(err).Msgf("cannot stop %q", c.Name)
				errCh <- err
			}
		}
		okCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return models.ErrShutdownTimeout
	case err := <-errCh:
		return err
	case <-okCh:
		return nil
	}
}
