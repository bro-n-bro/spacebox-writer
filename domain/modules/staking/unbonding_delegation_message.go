package staking

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
)

func UnbondingDelegationMessageHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.UnbondingDelegationMessage{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	coinBytes, err := jsoniter.Marshal(val.Coin)
	if err != nil {
		return err
	}

	val2 := storageModel.UnbondingDelegationMessage{
		CompletionTimestamp: val.CompletionTimestamp,
		Coin:                string(coinBytes),
		DelegatorAddress:    val.DelegatorAddress,
		ValidatorAddress:    val.ValidatorAddress,
		Height:              val.Height,
		TxHash:              val.TxHash,
	}

	var (
		db = ch.GetGormDB(ctx)
	)

	if err = db.Table("unbonding_delegation_message").Create(val2).Error; err != nil {
		return err
	}

	return nil
}
