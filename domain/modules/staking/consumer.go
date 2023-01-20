package staking

import (
	"context"

	"github.com/rs/zerolog"
	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/internal/configs"
)

type Consumer interface {
	subscribe(configs.Config, *clickhouse.Clickhouse, *zerolog.Logger) error
	handle(context.Context)
}
