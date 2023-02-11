package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

const (
	keyCMP     = "cmp"
	keyMetrics = "metrics"
)

type (
	Metrics struct {
		log *zerolog.Logger
		srv *http.Server

		cfg Config
	}
)

// New is a constructor for Metrics
func New(cfg Config, l zerolog.Logger) *Metrics {
	l = l.With().Str(keyCMP, keyMetrics).Logger()

	return &Metrics{
		log: &l,
		cfg: cfg,
	}
}

// Start is a method for starting metrics server
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

	return nil
}

// Stop is a method for stopping metrics server
func (m *Metrics) Stop(ctx context.Context) error {
	if !m.cfg.MetricsEnabled {
		return nil
	}

	return m.srv.Shutdown(ctx)
}
