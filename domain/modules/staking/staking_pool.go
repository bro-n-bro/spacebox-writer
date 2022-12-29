package staking

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"spacebox-writer/adapter/clickhouse"
)

func StakingPoolHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.StakingPool{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	var (
		db = ch.GetGormDB(ctx)
	)

	if err := db.Table("staking_pool").Create(val).Error; err != nil {
		return err
	}

	return nil
}
