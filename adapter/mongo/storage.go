package mongo

import (
	"context"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	msgStarted = "mongo client started"
)

type (
	Mongo struct {
		log        *zerolog.Logger
		cli        *mongo.Client
		collection *mongo.Collection
		cfg        Config
	}
)

// New is a constructor for Mongo
func New(cfg Config, l zerolog.Logger) *Mongo {
	l = l.With().Str("cmp", "mongo").Logger()
	return &Mongo{
		cfg: cfg,
		log: &l,
	}
}

// Start is a method for starting mongo client
func (s *Mongo) Start(ctx context.Context) error {
	opts := []*options.ClientOptions{
		options.Client().ApplyURI(s.cfg.URI),
		options.Client().SetMaxPoolSize(s.cfg.MaxPoolSize),
		options.Client().SetMaxConnecting(s.cfg.MaxConnecting),
		options.Client().SetAuth(options.Credential{
			Username: s.cfg.User,
			Password: s.cfg.Password,
		}),
	}

	client, err := mongo.Connect(ctx, opts...)
	if err != nil {
		return err
	}
	s.cli = client

	if err := s.Ping(ctx); err != nil {
		return err
	}

	collection := s.cli.Database("spacebox").Collection("handle_errors")
	s.collection = collection

	s.log.Info().Msg(msgStarted)

	return nil
}

// Stop is a method for stopping mongo client
func (s *Mongo) Stop(ctx context.Context) error {
	return s.cli.Disconnect(ctx)
}

// Ping is a method for checking connection to mongo
func (s *Mongo) Ping(ctx context.Context) error {
	return s.cli.Ping(ctx, nil)
}
