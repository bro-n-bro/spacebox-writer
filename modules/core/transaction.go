package core

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox-writer/internal/rep"
	"github.com/hexy-dev/spacebox/broker/model"
)

func TransactionHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.Transaction{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	return ch.Transaction(val)
}
