package staking

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	"spacebox-writer/internal/configs"

	"github.com/rs/zerolog"
)

type Consumer interface {
	subscribe(configs.Config, *clickhouse.Clickhouse, *zerolog.Logger) error
	handle(context.Context)
}
