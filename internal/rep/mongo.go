package rep

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/adapter/mongo/model"
)

type (
	// Mongo is a repository for mongo
	Mongo interface {
		HasBrokerMessage(ctx context.Context, id string) (bool, error)
		CreateBrokerMessage(ctx context.Context, msg *model.BrokerMessage) error
		DeleteBrokerMessage(ctx context.Context, id string) error
		UpdateBrokerMessage(ctx context.Context, msg *model.BrokerMessage) error
	}
)
