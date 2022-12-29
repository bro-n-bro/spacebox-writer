package mint

import (
	"context"
	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"spacebox-writer/adapter/clickhouse"
)

func AnnualProvisionHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.AnnualProvision{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	if err := ch.GetGormDB(ctx).Table("mint_params").Create(val).Error; err != nil {
		return err
	}

	return nil
}
