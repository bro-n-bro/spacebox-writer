package staking

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// DelegationMessageHandler is a handler for delegation message event
func DelegationMessageHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.DelegationMessage{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	return ch.DelegationMessage(val)
}
