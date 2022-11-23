package app

import (
	"spacebox-writer/domain/broker"
	"spacebox-writer/domain/graphql"
	"time"

	clhs "spacebox-writer/adapter/clickhouse"
)

type Config struct {
	Clickhouse   clhs.Config    `yaml:"clickhouse"`
	GraphQL      graphql.Config `yaml:"graphql"`
	Broker       broker.Config  `yaml:"broker"`
	StartTimeout time.Duration  `yaml:"start_timeout"`
	StopTimeout  time.Duration  `yaml:"stop_timeout"`
}
