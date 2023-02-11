package mint

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// MintParamsHandler is a handler for mint params event
func MintParamsHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.MintParams{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	return ch.MintParams(val)
}
