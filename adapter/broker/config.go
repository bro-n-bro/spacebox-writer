package broker

import "time"

// Config is a configuration for broker.
type Config struct {
	Address         string `env:"BROKER_SERVER"`
	GroupID         string `env:"GROUP_ID"`
	AutoOffsetReset string `env:"AUTO_OFFSET_RESET"`
	MaxRetries      int    `env:"MAX_RETRIES" envDefault:"5"`
	MetricsEnabled  bool   `env:"METRICS_ENABLED" envDefault:"false"`

	BatchBufferSize     int           `env:"BATCH_BUFFER_SIZE" envDefault:"200"`
	FlushBufferInterval time.Duration `env:"BATCH_FLUSH_BUFFER_INTERVAL" envDefault:"1m"`
}
