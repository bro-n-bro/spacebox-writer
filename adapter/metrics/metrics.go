package metrics

import (
	"context"
	"net/http"
	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/adapter/mongo"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

type Metrics struct {
	log   *zerolog.Logger
	srv   *http.Server
	ch    *clickhouse.Clickhouse
	mongo *mongo.Mongo

	stopScraping chan struct{}

	cfg Config
}

func New(cfg Config, l zerolog.Logger) *Metrics {
	l = l.With().Str("cmp", "metrics").Logger()

	return &Metrics{
		log:          &l,
		cfg:          cfg,
		stopScraping: make(chan struct{}),
	}
}

func (m *Metrics) Start(ctx context.Context) error {
	if !m.cfg.MetricsEnabled {
		return nil
	}

	m.srv = &http.Server{
		Addr:              ":" + m.cfg.Port,
		ReadHeaderTimeout: 1 * time.Second,
	}

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		if err := m.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			m.log.Fatal().Err(err).Msg("ListenAndServe error")
		}
	}()

	go m.startScraping()

	return nil
}

func (m *Metrics) Stop(ctx context.Context) error {
	if !m.cfg.MetricsEnabled {
		return nil
	}

	m.stopScraping <- struct{}{}

	return m.srv.Shutdown(ctx)
}