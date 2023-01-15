package broker

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"

	"github.com/hexy-dev/spacebox-writer/adapter/clickhouse"
	"github.com/hexy-dev/spacebox-writer/internal/rep"
)

type Broker struct {
	log       *zerolog.Logger
	pr        *kafka.Producer
	st        *clickhouse.Clickhouse
	m         rep.Mongo
	cfg       Config
	consumers []*kafka.Consumer
}

func New(cfg Config, st *clickhouse.Clickhouse, m rep.Mongo, log zerolog.Logger) *Broker {
	log = log.With().Str("cmp", "broker").Logger()
	return &Broker{
		log: &log,
		cfg: cfg,
		st:  st,
		m:   m,
	}
}

func (b *Broker) Start(_ context.Context) error {
	var err error
	b.pr, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": b.cfg.Address,
	})

	go func(drs chan kafka.Event) {
		for ev := range drs {
			m, ok := ev.(*kafka.Message)
			if !ok {
				continue
			}
			if err = m.TopicPartition.Error; err != nil {
				b.log.Error().Err(err).Msgf("Delivery error: %v", m.TopicPartition)
			}
		}
	}(b.pr.Events())

	return err
}

func (b *Broker) Stop(ctx context.Context) error {
	b.pr.Close()
	for _, consumer := range b.consumers {
		if err := consumer.Close(); err != nil {
			return err
		}
	}
	return nil
}
