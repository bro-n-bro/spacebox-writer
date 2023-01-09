package distribution

import (
	"context"

	"spacebox-writer/adapter/clickhouse"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func ValidatorCommissionHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.ValidatorCommission{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	if err := ch.GetGormDB(ctx).Table("validator_commission").Create(val).Error; err != nil {
		return err
	}

	return nil
}
