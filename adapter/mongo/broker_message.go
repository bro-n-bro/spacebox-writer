package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"spacebox-writer/adapter/mongo/model"
)

const (
	keyID = "_id"
)

func (s *Storage) HasBrokerMessage(ctx context.Context, id string) (r bool, err error) {
	msg := model.BrokerMessage{}

	if err = s.collection.FindOne(ctx, bson.D{{Key: keyID, Value: id}}).Decode(&msg); err == nil {
		return true, nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}

	return false, err
}

func (s *Storage) CreateBrokerMessage(ctx context.Context, msg *model.BrokerMessage) error {
	if _, err := s.collection.InsertOne(ctx, msg); err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateBrokerMessage(ctx context.Context, msg *model.BrokerMessage) error {
	filter := bson.D{{Key: keyID, Value: msg.ID}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "last_error_message", Value: msg.LastErrorMessage},
			{Key: "topic", Value: msg.Topic},
			{Key: "data", Value: msg.Data},
		}}}

	if _, err := s.collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (s *Storage) DeleteBrokerMessage(ctx context.Context, id string) error {
	filter := bson.D{{Key: keyID, Value: id}}
	if _, err := s.collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}
