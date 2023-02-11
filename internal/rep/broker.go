package rep

import (
	"context"
	"sync"
)

type (
	// Handler describes a handler for broker message
	Handler func(ctx context.Context, msg []byte, db Storage) error

	// Broker is a broker interface
	Broker interface {
		Subscribe(ctx context.Context, wg *sync.WaitGroup, topic string,
			handler func(ctx context.Context, msg []byte, db Storage) error) error
	}
)
