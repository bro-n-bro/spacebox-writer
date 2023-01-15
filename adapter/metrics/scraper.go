package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "spacebox_writer"
)

func (m *Metrics) startScraping() {
	m.log.Info().Msg("start metrics scraper")

	var (
		// last processed height
		heightMetric = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "last_processed_block_height",
			Help:      "Last processed block height",
		})

		lastHeight int64
		err        error
		ctx        = context.Background()
		ticker     = time.NewTicker(1 * time.Minute)
	)

	for {
		select {
		case <-m.stopScraping:
			m.log.Info().Msg("stop metrics scraper")
			return
		case <-ticker.C:
			// FIXME: it does not check handle errors

			if err = m.ch.GetGormDB(ctx).Select("height").Table("block").Order("height DESK").
				Limit(1).Scan(&lastHeight).Error; err != nil {
				m.log.Error().Err(err).Msg("get last block height from clickhouse error")
				continue
			}

			if lastHeight != 0 {
				heightMetric.Set(float64(lastHeight))
			}
		}
	}
}
