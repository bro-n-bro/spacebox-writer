package staking

import (
	"context"

	jsoniter "github.com/json-iterator/go"

	"github.com/hexy-dev/spacebox-writer/adapter/clickhouse"
	"github.com/hexy-dev/spacebox/broker/model"
)

func StakingParamsHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.StakingParams{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	return ch.StakingParams(val)
}
