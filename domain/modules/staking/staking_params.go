package staking

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
)

func StakingParamsHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.StakingParams{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	paramsBytes, err := jsoniter.Marshal(val.Params)
	if err != nil {
		return err
	}

	val2 := storageModel.StakingParams{
		Params: string(paramsBytes),
		Height: val.Height,
	}

	var (
		count int64
		db    = ch.GetGormDB(ctx)
	)

	if db.Table("staking_params").
		Where("height = ?", val.Height).
		Count(&count); count != 0 {

		return nil

	}

	if err = db.Table("staking_params").Create(val2).Error; err != nil {
		return errors.Wrap(err, "create staking_params error")
	}

	return nil
}
