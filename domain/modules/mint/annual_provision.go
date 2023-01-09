package mint

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"

	"spacebox-writer/adapter/clickhouse"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
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
