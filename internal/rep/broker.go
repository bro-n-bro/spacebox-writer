package rep

import (
	"context"
	"sync"

	"github.com/hexy-dev/spacebox-writer/adapter/clickhouse"
)

type (
	Handler func(ctx context.Context, msg []byte, db *clickhouse.Clickhouse) error

	Broker interface {
		Subscribe(ctx context.Context, wg *sync.WaitGroup, topic string,
			handler func(ctx context.Context, msg []byte, db *clickhouse.Clickhouse) error) error
	}
)
