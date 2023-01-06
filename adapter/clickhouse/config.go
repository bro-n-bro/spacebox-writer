package clickhouse

type Config struct {
	DSN            string `env:"CLICKHOUSE_DSN"`
	MigrationsPath string `env:"MIGRATIONS_PATH"`
	AutoMigrate    bool   `env:"AUTO_MIGRATE"`
}
