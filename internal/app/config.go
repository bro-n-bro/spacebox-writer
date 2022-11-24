package app

import (
	"spacebox-writer/domain/broker"
	"time"

	clhs "spacebox-writer/adapter/clickhouse"
)

type Config struct {
	Clickhouse   clhs.Config   `yaml:"clickhouse"`
	Broker       broker.Config `yaml:"broker"`
	StartTimeout time.Duration `yaml:"start_timeout"`
	StopTimeout  time.Duration `yaml:"stop_timeout"`
}
