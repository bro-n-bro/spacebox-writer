package staking

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox/broker/model"
	"spacebox-writer/adapter/clickhouse"
)

func StakingPoolHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.StakingPool{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	var (
		count int64
		db    = ch.GetGormDB(ctx)
	)

	if db.Table("staking_pool").
		Where("height = ?", val.Height).
		Count(&count); count != 0 {
		return nil

	}

	if err := db.Table("staking_pool").Create(val).Error; err != nil {
		return errors.Wrap(err, "create staking_pool error")
	}

	return nil
}
