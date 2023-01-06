package staking

//nolint:dupl

import (
	"context"

	"spacebox-writer/adapter/clickhouse"

	"github.com/hexy-dev/spacebox/broker/model"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func ValidatorStatusHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.ValidatorStatus{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	var (
		updates model.ValidatorStatus
		getVal  model.ValidatorStatus
		db      = ch.GetGormDB(ctx)
	)

	if err := db.Table("validator_status").
		Where("validator_address = ?", val.ValidatorAddress).
		First(&getVal).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = db.Table("validator_status").Create(val).Error; err != nil {
				return errors.Wrap(err, "create validator_status erorr")
			}
			return nil
		}
		return err
	}

	if val.Height > getVal.Height {

		if err := copier.Copy(&val, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}

		if err := db.Table("validator_status").
			Where("validator_address = ?", val.ValidatorAddress).
			Updates(&updates).Error; err != nil {

			return errors.Wrap(err, "update validator_status error")
		}
	}

	return nil
}
