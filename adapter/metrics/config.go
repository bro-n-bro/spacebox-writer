package metrics

// Config is a configuration for metrics.
type Config struct {
	Port           string `env:"METRICS_PORT" envDefault:"8080"`
	MetricsEnabled bool   `env:"METRICS_ENABLED" envDefault:"false"`
}
