package clickhouse

import "time"

// Config is a configuration for ClickHouse.
type Config struct {
	Addr                       string        `env:"CLICKHOUSE_ADDR" envDefault:"127.0.0.1:9000"`
	Database                   string        `env:"CLICKHOUSE_DATABASE" envDefault:"spacebox"`
	User                       string        `env:"CLICKHOUSE_USER" envDefault:"default"`
	Password                   string        `env:"CLICKHOUSE_PASSWORD"`
	MigrationsPath             string        `env:"MIGRATIONS_PATH"`
	BrokerServerForKafkaEngine string        `env:"BROKER_SERVER_FOR_KAFKA_ENGINE" envDefault:""`
	MaxIdleConns               int           `env:"CLICKHOUSE_MAX_IDLE_CONNS" envDefault:"20"`
	MaxOpenConns               int           `env:"CLICKHOUSE_MAX_OPEN_CONNS" envDefault:"25"`
	MaxExecutionTime           int           `env:"CLICKHOUSE_MAX_EXECUTION_TIME" envDefault:"60"`
	DialTimeout                time.Duration `env:"CLICKHOUSE_DIAL_TIMEOUT" envDefault:"10s"`
	AutoMigrate                bool          `env:"AUTO_MIGRATE"`
}
