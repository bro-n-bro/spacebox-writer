package broker

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog"

	"github.com/bro-n-bro/spacebox-writer/adapter/clickhouse"
	"github.com/bro-n-bro/spacebox-writer/internal/rep"
)

const (
	msgDeliveryError        = "delivery error: %v"
	msgFlushedKafkaMessages = "flushed kafka messages. Outstanding events still un-flushed: %d"
	msgKafkaLocalQueueFull  = "kafka local queue full error - Going to Flush then retry"
)

type (
	Broker struct {
		m         rep.Mongo
		metrics   *metrics
		log       *zerolog.Logger
		pr        *kafka.Producer
		st        *clickhouse.Clickhouse
		consumers []*kafka.Consumer
		cfg       Config
	}

	metrics struct {
		histogram            *prometheus.HistogramVec
		errorsCounter        *prometheus.CounterVec
		limitExceededCounter *prometheus.CounterVec
	}
)

// New creates new broker instance.
func New(cfg Config, st *clickhouse.Clickhouse, m rep.Mongo, log zerolog.Logger) *Broker {
	log = log.With().Str("cmp", "broker").Logger()

	b := &Broker{
		log: &log,
		cfg: cfg,
		st:  st,
		m:   m,
	}

	if b.cfg.MetricsEnabled {
		b.metrics = &metrics{
			histogram: promauto.NewHistogramVec(prometheus.HistogramOpts{
				Namespace: "spacebox_writer",
				Name:      "process_duration",
				Help:      "Duration of handling kafka messages.",
			}, []string{keyTopic}),
			errorsCounter: promauto.NewCounterVec(prometheus.CounterOpts{
				Namespace: "spacebox_writer",
				Name:      "fails_total",
				Help:      "Count of handling errors for each topic.",
			}, []string{keyTopic}),
			limitExceededCounter: promauto.NewCounterVec(prometheus.CounterOpts{
				Namespace: "spacebox_writer",
				Name:      "exceeded_limit_total",
				Help:      "Count of limits exceeded for each topic.",
			}, []string{keyTopic}),
		}
	}

	return b
}

// Start starts broker.
func (b *Broker) Start(_ context.Context) (err error) {
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
				b.log.Error().Err(err).Msgf(msgDeliveryError, m.TopicPartition)
			}
		}
	}(b.pr.Events())

	return err
}

// Stop stops broker.
func (b *Broker) Stop(ctx context.Context) error {
	b.pr.Close()
	for _, consumer := range b.consumers {
		if err := consumer.Close(); err != nil {
			return err
		}
	}
	return nil
}
