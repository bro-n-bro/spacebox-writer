package auth

import (
	"context"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/hexy-dev/spacebox/broker/model"
	"spacebox-writer/adapter/clickhouse"
)

func AccountHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.Account{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	var (
		db      = ch.GetGormDB(ctx)
		updates model.Account
		getVal  model.Account
	)

	if err := db.Table("account").
		Where("address = ?", val.Address).
		First(&getVal).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = db.Table("account").Create(val).Error; err != nil {
				return errors.Wrap(err, "error of create")
			}
			return nil
		}
		return errors.Wrap(err, "error of database: %w")
	}

	if val.Height < getVal.Height {
		if err := copier.Copy(&val, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}

		if err := db.Table("account").
			Where("address = ?", val.Height).
			Updates(&updates).Error; err != nil {

			return errors.Wrap(err, "error of update")
		}
	}

	return nil
}
