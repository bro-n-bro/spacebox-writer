package app

import (
	"time"

	"github.com/bro-n-bro/spacebox-writer/adapter/broker"
	"github.com/bro-n-bro/spacebox-writer/adapter/clickhouse"
	"github.com/bro-n-bro/spacebox-writer/adapter/metrics"
	"github.com/bro-n-bro/spacebox-writer/adapter/mongo"
	"github.com/bro-n-bro/spacebox-writer/modules"
)

type Config struct {
	LogLevel     string `env:"LOG_LEVEL"`
	Metrics      metrics.Config
	Broker       broker.Config
	Modules      modules.Config
	Mongo        mongo.Config
	Clickhouse   clickhouse.Config
	StartTimeout time.Duration `env:"START_TIMEOUT"`
	StopTimeout  time.Duration `env:"STOP_TIMEOUT"`
}
