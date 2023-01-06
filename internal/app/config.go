package app

import (
	"spacebox-writer/adapter/mongo"
	"spacebox-writer/domain/modules"
	"time"

	"spacebox-writer/adapter/broker"
	"spacebox-writer/adapter/clickhouse"
)

type Config struct {
	Clickhouse   clickhouse.Config
	Broker       broker.Config
	Mongo        mongo.Config
	Modules      modules.Config
	StartTimeout time.Duration `env:"START_TIMEOUT"`
	StopTimeout  time.Duration `env:"STOP_TIMEOUT"`
	LogLevel     string        `env:"LOG_LEVEL"`
}
