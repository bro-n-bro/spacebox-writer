package app

import (
	"time"

	"spacebox-writer/adapter/broker"
	"spacebox-writer/domain/modules"

	clhs "spacebox-writer/adapter/clickhouse"
)

type Config struct {
	Clickhouse   clhs.Config    `yaml:"clickhouse"`
	Broker       broker.Config  `yaml:"broker"`
	StartTimeout time.Duration  `yaml:"start_timeout"`
	StopTimeout  time.Duration  `yaml:"stop_timeout"`
	Modules      modules.Config `yaml:"modules"`
}
