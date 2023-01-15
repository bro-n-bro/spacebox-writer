package distribution

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox-writer/adapter/clickhouse"
	"github.com/hexy-dev/spacebox/broker/model"
)

func ValidatorCommissionHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.ValidatorCommission{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	return ch.ValidatorCommission(val)
}
