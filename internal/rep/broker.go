package rep

import (
	"context"
	"sync"
)

type (
	Handler func(ctx context.Context, msg []byte, db Storage) error

	Broker interface {
		Subscribe(ctx context.Context, wg *sync.WaitGroup, topic string,
			handler func(ctx context.Context, msg []byte, db Storage) error) error
	}
)
