package mongo

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/adapter/mongo/model"
)

// CreateBrokerMessage is a method for creating a broker message in the database
func (s *Mongo) CreateBrokerMessage(ctx context.Context, msg *model.BrokerMessage) error {
	if _, err := s.collection.InsertOne(ctx, msg); err != nil {
		return err
	}

	return nil
}
