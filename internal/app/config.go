package app

import (
	"time"

	"github.com/hexy-dev/spacebox-writer/adapter/broker"
	"github.com/hexy-dev/spacebox-writer/adapter/clickhouse"
	"github.com/hexy-dev/spacebox-writer/adapter/mongo"
	"github.com/hexy-dev/spacebox-writer/modules"
)

type Config struct {
	LogLevel     string `env:"LOG_LEVEL"`
	Broker       broker.Config
	Modules      modules.Config
	Mongo        mongo.Config
	Clickhouse   clickhouse.Config
	StartTimeout time.Duration `env:"START_TIMEOUT"`
	StopTimeout  time.Duration `env:"STOP_TIMEOUT"`
}
