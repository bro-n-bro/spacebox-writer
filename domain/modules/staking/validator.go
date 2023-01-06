package staking

import (
	"context"

	"github.com/hexy-dev/spacebox/broker/model"

	"spacebox-writer/adapter/clickhouse"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func ValidatorHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.Validator{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	var (
		count int64
		db    = ch.GetGormDB(ctx)
	)

	if db.Table("validator").
		Where("consensus_address = ?", val.ConsensusAddress).
		Count(&count); count != 0 {

		return nil
	}

	if err := db.Table("validator").Create(val).Error; err != nil {
		return errors.Wrap(err, "create validator error")
	}

	return nil
}
