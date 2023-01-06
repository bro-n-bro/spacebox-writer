package mongo

import (
	"context"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	log        *zerolog.Logger
	cli        *mongo.Client
	collection *mongo.Collection
	cfg        Config
}

func New(cfg Config, l zerolog.Logger) *Storage {
	l = l.With().Str("cmp", "mongo").Logger()
	return &Storage{
		cfg: cfg,
		log: &l,
	}
}

func (s *Storage) Start(ctx context.Context) error {
	opts := []*options.ClientOptions{
		options.Client().ApplyURI(s.cfg.URI),
		options.Client().SetMaxPoolSize(s.cfg.MaxPoolSize),
		options.Client().SetMaxConnecting(s.cfg.MaxConnecting),
		options.Client().SetAuth(options.Credential{
			AuthMechanism:           "",
			AuthMechanismProperties: nil,
			AuthSource:              "",
			Username:                s.cfg.User,
			Password:                s.cfg.Password,
			PasswordSet:             false,
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

	s.log.Info().Msg("storage started")

	return nil
}

func (s *Storage) Stop(ctx context.Context) error {
	return s.cli.Disconnect(ctx)
}

func (s *Storage) Ping(ctx context.Context) error {
	return s.cli.Ping(ctx, nil)
}
