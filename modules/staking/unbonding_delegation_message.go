package staking

import (
	"context"

	jsoniter "github.com/json-iterator/go"

	"github.com/hexy-dev/spacebox-writer/adapter/clickhouse"
	"github.com/hexy-dev/spacebox/broker/model"
)

func UnbondingDelegationMessageHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.UnbondingDelegationMessage{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	return ch.UnbondingDelegationMessage(val)
}
