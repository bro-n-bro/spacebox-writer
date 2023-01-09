package staking

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
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

	if err = ch.GetGormDB(ctx).Table("delegation_message").Create(storageModel.DelegationMessage{
		OperatorAddress:  val.OperatorAddress,
		DelegatorAddress: val.DelegatorAddress,
		Coin:             string(coinBytes),
		Height:           val.Height,
		TxHash:           val.TxHash,
	}).Error; err != nil {
		return err
	}

	return nil
}
