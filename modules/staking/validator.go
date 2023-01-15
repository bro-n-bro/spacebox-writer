package staking

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"spacebox-writer/adapter/clickhouse"

	"github.com/hexy-dev/spacebox/broker/model"
)

func ValidatorHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.Validator{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	return ch.Validator(val)
}
