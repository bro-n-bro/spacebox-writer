package clickhouse

type Config struct {
	DSN            string `yaml:"dsn"`
	MigrationsPath string `yaml:"migrations_path"`
	AutoMigrate    bool   `yaml:"auto_migrate"`
}
