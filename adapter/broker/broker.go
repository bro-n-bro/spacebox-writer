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
	p         *kafka.Producer
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
			if err := m.TopicPartition.Error; err != nil {
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

func (b *Broker) produce(topic string, data []byte, headers []kafka.Header) error {
	err := b.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          data,
		Headers:        headers,
	}, nil)

	if kafkaError, ok := err.(kafka.Error); ok && kafkaError.Code() == kafka.ErrQueueFull {
		b.log.Info().Str("topic", topic).Msg("Kafka local queue full error - Going to Flush then retry...")
		flushedMessages := b.p.Flush(30 * 1000)
		b.log.Info().Str("topic", topic).
			Msgf("Flushed kafka messages. Outstanding events still un-flushed: %d", flushedMessages)
		return b.produce(topic, data, headers)
	}

	return nil
}

func bytesToInt(bytes []byte) int {
	var value int
	for _, bt := range bytes {
		value = value*10 + int(bt-48)
	}
	return value
}
