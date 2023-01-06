package distribution

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func DistributionParamsHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.DistributionParams{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	paramsBytes, err := jsoniter.Marshal(val.Params)
	if err != nil {
		return err
	}

	val2 := storageModel.DistributionParams{
		Params: string(paramsBytes),
		Height: val.Height,
	}

	if err = ch.GetGormDB(ctx).Table("distribution_params").Create(val2).Error; err != nil {
		return err
	}

	return nil
}
