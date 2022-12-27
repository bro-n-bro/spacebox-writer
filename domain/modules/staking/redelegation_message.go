package staking

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"

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
		count int64
		db    = ch.GetGormDB(ctx)
	)

	if db.Table("redelegation_message").
		Where("height = ? AND src_validator_address = ? AND delegator_address = ? AND dst_validator_address = ?",
			val2.Height,
			val2.SrcValidatorAddress,
			val2.DelegatorAddress,
			val2.DstValidatorAddress,
		).Count(&count); count != 0 {
		return nil

	}

	if err = db.Table("redelegation_message").Create(val2).Error; err != nil {
		return errors.Wrap(err, "create redelegation_message error")
	}

	return nil
}
