package distribution

import (
	"context"
	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
)

func DelegationRewardHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.DelegationReward{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	paramsBytes, err := jsoniter.Marshal(val.Coins)
	if err != nil {
		return err
	}

	val2 := storageModel.DelegationReward{
		Coins:            string(paramsBytes),
		DelegatorAddress: val.DelegatorAddress,
		WithdrawAddress:  val.WithdrawAddress,
		OperatorAddress:  val.OperatorAddress,
		Height:           val.Height,
	}

	ch.GetGormDB(ctx).Table("delegation_reward").Create(val2)

	return nil

}
