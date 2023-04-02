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

type (
	Broker struct {
		m         rep.Mongo
		metrics   *metrics
		log       *zerolog.Logger
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
func (b *Broker) Start(_ context.Context) error {
	return nil
}

// Stop stops broker.
func (b *Broker) Stop(ctx context.Context) error {
	for _, consumer := range b.consumers {
		if err := consumer.Close(); err != nil {
			return err
		}
	}
	return nil
}
