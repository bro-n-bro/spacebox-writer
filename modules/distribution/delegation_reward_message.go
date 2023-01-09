package distribution

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func DelegationRewardMessageHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.DelegationRewardMessage{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	paramsBytes, err := jsoniter.Marshal(val.Coins)
	if err != nil {
		return err
	}

	if err = ch.GetGormDB(ctx).Table("delegation_reward_message").Create(storageModel.DelegationRewardMessage{
		Coins:            string(paramsBytes),
		Height:           val.Height,
		DelegatorAddress: val.DelegatorAddress,
		ValidatorAddress: val.ValidatorAddress,
		TxHash:           val.TxHash,
	}).Error; err != nil {
		return err
	}

	return nil
}
