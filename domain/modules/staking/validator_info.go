package staking

import (
	"context"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/hexy-dev/spacebox/broker/model"
	"spacebox-writer/adapter/clickhouse"
)

func ValidatorInfoHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.ValidatorInfo{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	var (
		db      = ch.GetGormDB(ctx)
		updates model.ValidatorInfo
		getVal  model.ValidatorInfo
	)

	if err := db.Table("validator_info").
		Where("consensus_address = ?", val.ConsensusAddress).
		First(&getVal).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = db.Table("validator_info").Create(val).Error; err != nil {
				return errors.Wrap(err, "create validator_info error")
			}
			return nil
		}
		return err
	}

	if val.Height < getVal.Height {
		if err := copier.Copy(&val, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}

		if err := db.Table("validator_info").
			Where("consensus_address = ?", val.ConsensusAddress).
			Updates(&updates).Error; err != nil {

			return errors.Wrap(err, "update validator_info error")
		}
	}

	return nil
}
