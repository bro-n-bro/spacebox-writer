package graphql

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/graph/generated"
	resolvers "spacebox-writer/graph/resolvers"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type GraphQL struct {
	log *zerolog.Logger
	db  *gorm.DB
	cfg Config
}

func New(cfg Config, db *gorm.DB) *GraphQL {
	lg := zerolog.New(os.Stderr).Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().
		Str("cmp", "graphql").Logger()
	return &GraphQL{
		log: &lg,
		cfg: cfg,
		db:  db,
	}
}

func (gql *GraphQL) Start(context.Context) error {
	errCh := make(chan error)
	gql.log.Debug().Msgf("start listening %q", gql.cfg.Address)
	go func() {

		srv := handler.NewDefaultServer(
			generated.NewExecutableSchema(
				generated.Config{
					Resolvers: &resolvers.Resolver{},
				},
			),
		)

		customCtx := &clickhouse.CustomContext{
			Database: gql.db,
		}

		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		http.Handle("/query", clickhouse.CreateContext(customCtx, srv))
		if err := http.ListenAndServe(gql.cfg.Address, nil); err != nil {
			errCh <- errors.Wrap(err, "cannot listen and serve")

		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-time.After(gql.cfg.StartTimeout):
		return nil
	}

}

func (gql *GraphQL) Stop(ctx context.Context) error { return nil }
