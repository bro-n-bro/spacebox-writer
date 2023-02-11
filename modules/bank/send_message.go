package bank

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// SendMessageHandler is a handler for send message event
func SendMessageHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.SendMessage{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	return ch.SendMessage(val)
}
