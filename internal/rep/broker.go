package rep

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
)

type Handler func(ctx context.Context, msg []byte, db *clickhouse.Clickhouse) error

type Broker interface {
	Subscribe(ctx context.Context, topic string,
		handler func(ctx context.Context, msg []byte, db *clickhouse.Clickhouse) error) error
}
