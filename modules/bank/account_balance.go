package bank

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox-writer/adapter/clickhouse"
	"github.com/hexy-dev/spacebox/broker/model"
)

func AccountBalanceHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.AccountBalance{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	return ch.AccountBalance(val)
}
