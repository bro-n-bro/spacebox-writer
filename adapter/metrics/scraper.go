package metrics

import (
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
		ticker     = time.NewTicker(1 * time.Minute)
	)

	for {
		select {
		case <-m.stopScraping:
			m.log.Info().Msg("stop metrics scraper")
			return
		case <-ticker.C:
			// FIXME: it does not check handle errors

			lastHeight, err = m.ch.LatestBlockHeight()
			if err != nil {
				m.log.Error().Err(err).Msg("get LatestBlockHeight from clickhouse error")
				continue
			}

			if lastHeight != 0 {
				heightMetric.Set(float64(lastHeight))
			}
		}
	}
}
