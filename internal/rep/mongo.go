package rep

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/adapter/mongo/model"
)

type (
	// Mongo is a repository for mongo
	Mongo interface {
		CreateBrokerMessage(ctx context.Context, msg *model.BrokerMessage) error
	}
)
