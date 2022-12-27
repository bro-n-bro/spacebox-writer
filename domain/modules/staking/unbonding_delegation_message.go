package staking

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/hexy-dev/spacebox/broker/model"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
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
		getVal storageModel.UnbondingDelegationMessage
		db     = ch.GetGormDB(ctx)
	)

	if err := db.Table("unbonding_delegation_message").
		Where("validator_address = ? AND delegator_address = ? AND height = ?",
			val2.ValidatorAddress,
			val2.DelegatorAddress,
			val2.Height,
		).First(&getVal).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = db.Table("unbonding_delegation_message").Create(val2).Error; err != nil {
				return errors.Wrap(err, "create unbonding_delegation_message error")
			}
			return nil
		}
		return err
	}

	return nil
}
