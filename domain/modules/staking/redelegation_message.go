package staking

import (
	"context"
	"encoding/json"
	"github.com/hexy-dev/spacebox/broker/model"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
)

func RedelegationMessageHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.RedelegationMessage{}
	if err := json.Unmarshal(msg, &val); err != nil {
		return err
	}

	coinBytes, err := json.Marshal(val.Coin)
	if err != nil {
		return err
	}

	val2 := storageModel.RedelegationMessage{
		CompletionTime:      val.CompletionTime,
		Coin:                string(coinBytes),
		DelegatorAddress:    val.DelegatorAddress,
		SrcValidatorAddress: val.SrcValidatorAddress,
		DstValidatorAddress: val.DstValidatorAddress,
		Height:              val.Height,
		TxHash:              val.TxHash,
	}

	var (
		db = ch.GetGormDB(ctx)
	)

	if err = db.Table("redelegation_message").Create(val2).Error; err != nil {
		return err
	}

	return nil
}
