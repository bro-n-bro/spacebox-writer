package staking

import (
	"context"

	"spacebox-writer/adapter/clickhouse"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
)

func StakingPoolHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.StakingPool{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	if err := ch.GetGormDB(ctx).Table("staking_pool").Create(val).Error; err != nil {
		return err
	}

	return nil
}
