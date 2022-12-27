package staking

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
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

	val2 := storageModel.Redelegation{
		CompletionTime:      val.CompletionTime,
		Coin:                string(coinBytes),
		DelegatorAddress:    val.DelegatorAddress,
		SrcValidatorAddress: val.SrcValidatorAddress,
		DstValidatorAddress: val.DstValidatorAddress,
		Height:              val.Height,
	}

	ch.GetGormDB(ctx).Table("redelegation").Create(val2)

	// v.db.SaveValidator() // interface implementation in adapter
	return nil
}
