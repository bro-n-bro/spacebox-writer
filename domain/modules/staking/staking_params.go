package staking

import (
	"context"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

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
		stVal   storageModel.StakingParams
		updates storageModel.StakingParams
		db      = ch.GetGormDB(ctx)
	)

	err = db.Table("staking_params").
		Where("height = ?", val.Height).
		First(&stVal).Error

	if !errors.Is(gorm.ErrRecordNotFound, err) {
		return err
	} else if val2.Params != stVal.Params {
		if err = copier.Copy(&val2, &updates); err != nil {
			return err
		}
		if err = db.Table("staking_params").
			Where("height = ?", val2.Height).
			Updates(&updates).Error; err != nil {
			return err
		}
	}

	return nil
}
