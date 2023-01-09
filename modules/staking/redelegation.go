package staking

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func RedelegationHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.Redelegation{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	coinBytes, err := jsoniter.Marshal(val.Coin)
	if err != nil {
		return errors.Wrap(err, "marshall error")
	}

	if err = ch.GetGormDB(ctx).Table("redelegation").Create(storageModel.Redelegation{
		CompletionTime:      val.CompletionTime,
		Coin:                string(coinBytes),
		DelegatorAddress:    val.DelegatorAddress,
		SrcValidatorAddress: val.SrcValidatorAddress,
		DstValidatorAddress: val.DstValidatorAddress,
		Height:              val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}
