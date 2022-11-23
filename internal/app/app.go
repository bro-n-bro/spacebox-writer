package app

import (
	"context"
	"os"
	"spacebox-writer/domain/broker"

	clhs "spacebox-writer/adapter/clickhouse"
	"spacebox-writer/consts"
	gql "spacebox-writer/domain/graphql"
	"spacebox-writer/internal/rep"
	"spacebox-writer/models"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type cmp struct {
	Service rep.Lifecycle
	Name    string
}

type App struct {
	log  *zerolog.Logger
	cmps []cmp
	cfg  Config
}

func New(cfg Config) *App {
	l := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "app").Logger()
	return &App{
		log:  &l,
		cfg:  cfg,
		cmps: []cmp{},
	}
}

func (a *App) Start(ctx context.Context) error {
	a.log.Info().Msg("starting app")

	clickhouse := clhs.New(a.cfg.Clickhouse)
	gormDB := clickhouse.GetGormDB(ctx)
	graphql := gql.New(a.cfg.GraphQL, gormDB)
	brk := broker.New(a.cfg.Broker)

	a.cmps = append(
		a.cmps,
		cmp{clickhouse, "clickhouse"},
		cmp{graphql, "graphql"},
		cmp{brk, "broker"},
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

	errCh := make(chan error)
	go func() {
		gr, ctx := errgroup.WithContext(ctx)
		var c cmp
		for i := len(a.cmps) - 1; i >= 0; i-- {
			c = a.cmps[i]
			a.log.Info().Msgf("stopping %q...", c.Name)
			if err := c.Service.Stop(ctx); err != nil {
				a.log.Error().Err(err).Msgf("cannot stop %q", c.Name)
			}
		}
		errCh <- gr.Wait()
	}()

	select {
	case <-ctx.Done():
		return models.ErrShutdownTimeout
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	}
}
