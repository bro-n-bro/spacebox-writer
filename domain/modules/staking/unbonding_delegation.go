package staking

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
)

func UnbondingDelegationHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.UnbondingDelegation{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	coinBytes, err := jsoniter.Marshal(val.Coin)
	if err != nil {
		return err
	}

	val2 := storageModel.UnbondingDelegation{
		CompletionTimestamp: val.CompletionTimestamp,
		Coin:                string(coinBytes),
		DelegatorAddress:    val.DelegatorAddress,
		ValidatorAddress:    val.ValidatorAddress,
		Height:              val.Height,
	}

	var (
		db = ch.GetGormDB(ctx)
	)

	if err = db.Table("unbonding_delegation").Create(val2).Error; err != nil {
		return err
	}

	return nil
}
