package broker

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rs/zerolog"

	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/internal/rep"
)

type Broker struct {
	log       *zerolog.Logger
	clh       *clickhouse.Clickhouse
	p         *kafka.Producer
	m         rep.Mongo
	cfg       Config
	consumers []*kafka.Consumer
}

func New(cfg Config, clickhouse *clickhouse.Clickhouse, m rep.Mongo, l zerolog.Logger) *Broker {
	l = l.With().Str("cmp", "broker").Logger()
	return &Broker{
		log: &l,
		cfg: cfg,
		clh: clickhouse,
		m:   m,
	}
}

func (b *Broker) Start(_ context.Context) error {
	var err error
	b.p, err = kafka.NewProducer(&kafka.ConfigMap{
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
	}(b.p.Events())

	return err
}

func (b *Broker) Stop(ctx context.Context) error {
	b.p.Close()
	for _, consumer := range b.consumers {
		if err := consumer.Close(); err != nil {
			return err
		}
	}
	return nil
}
