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

func DelegationMessageHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.DelegationMessage{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	coinBytes, err := jsoniter.Marshal(val.Coin)
	if err != nil {
		return errors.Wrap(err, "marshall error")
	}

	val2 := storageModel.DelegationMessage{
		OperatorAddress:  val.OperatorAddress,
		DelegatorAddress: val.DelegatorAddress,
		Coin:             string(coinBytes),
		Height:           val.Height,
		TxHash:           val.TxHash,
	}

	var (
		getVal storageModel.DelegationMessage
		db     = ch.GetGormDB(ctx)
	)

	if err := db.Table("delegation_message").
		Where("operator_address = ? AND delegator_address = ? AND height = ?",
			val2.OperatorAddress,
			val2.DelegatorAddress,
			val2.Height,
		).First(&getVal).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = db.Table("delegation_message").Create(val2).Error; err != nil {
				return errors.Wrap(err, "create delegation_message error")
			}
			return nil
		}
		return err
	}

	return nil
}
