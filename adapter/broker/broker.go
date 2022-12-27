package broker

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"
	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/internal/configs"
)

type Broker struct {
	log       *zerolog.Logger
	clh       *clickhouse.Clickhouse
	cfg       configs.Config
	consumers []*kafka.Consumer
}

func New(cfg configs.Config, clickhouse *clickhouse.Clickhouse, l zerolog.Logger) *Broker {
	l = l.With().Str("cmp", "broker").Logger()
	return &Broker{
		log: &l,
		cfg: cfg,
		clh: clickhouse,
	}
}

func (b *Broker) Subscribe(
	ctx context.Context,
	topic string,
	handler func(ctx context.Context, msg []byte, db *clickhouse.Clickhouse) error,
) error {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        b.cfg.Address,
		"group.id":                 b.cfg.GroupID,
		"auto.offset.reset":        b.cfg.AutoOffsetReset,
		"allow.auto.create.topics": true,
	})

	if err != nil {
		return err
	}

	if err = consumer.Subscribe(topic, nil); err != nil {
		return err
	}
	b.consumers = append(b.consumers, consumer)

	go func() {
		for {
			select {
			case <-ctx.Done():
				b.log.Info().Str("topic", topic).Msg("stop read messages from topic")
				return
			default:
			}

			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				b.log.
					Error().
					Err(err).
					Str("msg", string(msg.Value)).
					Msg("read message error")

				continue
			} else {
				b.log.Debug().Msgf("[%v]: %s", msg.String(), msg.Value)
			}

			if err = handler(ctx, msg.Value, b.clh); err != nil {
				b.log.
					Error().
					Err(err).
					Str("topic", topic).
					Str("msg", string(msg.Value)).
					Msg("handle message error")

				continue
			}

			_, err = consumer.CommitMessage(msg)
			if err != nil {
				b.log.
					Error().
					Err(err).
					Str("topic", topic).
					Msg("commit message error")
			}
		}
	}()

	return nil
}

func (b *Broker) Stop(ctx context.Context) error {
	for _, consumer := range b.consumers {
		if err := consumer.Close(); err != nil {
			return err
		}
	}
	return nil
}
