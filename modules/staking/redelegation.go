package staking

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// RedelegationHandler is a handler for redelegation event
func RedelegationHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.Redelegation{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	return ch.Redelegation(val)
}
