package staking

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	"encoding/json"
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

	if err = ch.GetGormDB(ctx).Table("redelegation_message").Create(storageModel.RedelegationMessage{
		CompletionTime:      val.CompletionTime,
		Coin:                string(coinBytes),
		DelegatorAddress:    val.DelegatorAddress,
		SrcValidatorAddress: val.SrcValidatorAddress,
		DstValidatorAddress: val.DstValidatorAddress,
		Height:              val.Height,
		TxHash:              val.TxHash,
	}).Error; err != nil {
		return err
	}

	return nil
}
