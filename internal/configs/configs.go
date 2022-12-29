package configs

import (
	"time"
)

type Config struct {
	DSN             string        `env:"CLICKHOUSE_DSN"`
	MigrationsPath  string        `env:"MIGRATIONS_PATH"`
	Address         string        `env:"BROKER_SERVER"`
	GroupID         string        `env:"GROUP_ID"`
	AutoOffsetReset string        `env:"AUTO_OFFSET_RESET"`
	Modules         []string      `env:"MODULES"`
	StartTimeout    time.Duration `env:"START_TIMEOUT"`
	StopTimeout     time.Duration `env:"STOP_TIMEOUT"`
	AutoMigrate     bool          `env:"AUTO_MIGRATE"`
	LogLevel        string        `env:"LOG_LEVEL"`
}
