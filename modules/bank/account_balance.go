package bank

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// AccountBalanceHandler is a handler for account balance event
func AccountBalanceHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.AccountBalance{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	return ch.AccountBalance(val)
}
