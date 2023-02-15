package mongo

import (
	"context"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bro-n-bro/spacebox-writer/adapter/mongo/model"
)

const (
	keyID = "_id"
)

// HasBrokerMessage is a method for checking if a broker message exists in the database
func (s *Mongo) HasBrokerMessage(ctx context.Context, id string) (r bool, err error) {
	msg := model.BrokerMessage{}

	if err = s.collection.FindOne(ctx, bson.D{{Key: keyID, Value: id}}).Decode(&msg); err == nil {
		return true, nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}

	return false, err
}

// CreateBrokerMessage is a method for creating a broker message in the database
func (s *Mongo) CreateBrokerMessage(ctx context.Context, msg *model.BrokerMessage) error {
	if _, err := s.collection.InsertOne(ctx, msg); err != nil {
		return err
	}

	return nil
}

// UpdateBrokerMessage is a method for updating a broker message in the database
func (s *Mongo) UpdateBrokerMessage(ctx context.Context, msg *model.BrokerMessage) error {
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

// GetMessagesByTopic is a method for getting broker messages by topic from the database
func (s *Mongo) GetMessagesByTopic(ctx context.Context, topic string) ([]*model.BrokerMessage, error) {
	filter := bson.D{{Key: "topic", Value: topic}}
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	messages := make([]*model.BrokerMessage, 0)
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// DeleteBrokerMessage is a method for deleting a broker message from the database
func (s *Mongo) DeleteBrokerMessage(ctx context.Context, id string) error {
	filter := bson.D{{Key: keyID, Value: id}}
	if _, err := s.collection.DeleteOne(ctx, filter); err != nil {
		return err
	}

	return nil
}

// DeleteBrokerMessages is a method for deleting a broker message from the database
func (s *Mongo) DeleteBrokerMessages(ctx context.Context, ids []string) error {
	filter := bson.D{{Key: keyID, Value: bson.D{{Key: "$in", Value: ids}}}}
	if _, err := s.collection.DeleteMany(ctx, filter); err != nil {
		return err
	}

	return nil
}
