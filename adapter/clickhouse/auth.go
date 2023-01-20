package clickhouse

import (
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/bro-n-bro/spacebox/broker/model"
)

func (ch *Clickhouse) Account(val model.Account) (err error) {
	var (
		updates model.Account
		getVal  model.Account
	)

	if err = ch.gorm.Table("account").
		Where("address = ?", val.Address).
		First(&getVal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table("account").Create(val).Error; err != nil {
				return errors.Wrap(err, "error of create")
			}
			return nil
		}
		return errors.Wrap(err, "error of database: %w")
	}

	if val.Height < getVal.Height {
		if err = copier.Copy(&val, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}

		if err = ch.gorm.Table("account").Where("address = ?", val.Height).Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "error of update")
		}
	}

	return nil
}
