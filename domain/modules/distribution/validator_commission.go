package distribution

import (
	"context"
	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"spacebox-writer/adapter/clickhouse"
)

func ValidatorCommissionHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.ValidatorCommission{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	ch.GetGormDB(ctx).Table("validator_commission").Create(val)

	return nil
}
