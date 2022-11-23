package rep

import "context"

// Lifecycle indicates that struct can be run as an application component
type Lifecycle interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
