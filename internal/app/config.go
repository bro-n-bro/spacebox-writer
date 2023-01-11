package app

import (
	"time"

	"spacebox-writer/adapter/metrics"
	"spacebox-writer/adapter/mongo"
	"spacebox-writer/modules"

	"spacebox-writer/adapter/broker"
	"spacebox-writer/adapter/clickhouse"
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
