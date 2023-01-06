package rep

import (
	"context"

	"spacebox-writer/adapter/mongo/model"
)

type Mongo interface {
	HasBrokerMessage(ctx context.Context, id string) (bool, error)
	CreateBrokerMessage(ctx context.Context, msg *model.BrokerMessage) error
	DeleteBrokerMessage(ctx context.Context, id string) error
	UpdateBrokerMessage(ctx context.Context, msg *model.BrokerMessage) error
}
