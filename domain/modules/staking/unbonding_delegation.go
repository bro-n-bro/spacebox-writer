package staking

import (
	"context"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/hexy-dev/spacebox/broker/model"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
)

func UnbondingDelegationHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.UnbondingDelegation{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	coinBytes, err := jsoniter.Marshal(val.Coin)
	if err != nil {
		return err
	}

	val2 := storageModel.UnbondingDelegation{
		CompletionTimestamp: val.CompletionTimestamp,
		Coin:                string(coinBytes),
		DelegatorAddress:    val.DelegatorAddress,
		ValidatorAddress:    val.ValidatorAddress,
		Height:              val.Height,
	}

	// if validator_address + delegator_address
	// not exists - create
	// if new height greater than height
	// in DB - update height

	var (
		updates storageModel.UnbondingDelegation
		getVal  storageModel.UnbondingDelegation
		db      = ch.GetGormDB(ctx)
	)

	if err := db.Table("unbonding_delegation").
		Where("validator_address = ? AND delegator_address = ?",
			val2.ValidatorAddress,
			val2.DelegatorAddress,
		).First(&getVal).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = db.Table("unbonding_delegation").Create(val2).Error; err != nil {
				return errors.Wrap(err, "create unbonding_delegation error")
			}
			return nil
		}
		return err
	}

	if val2.Height > getVal.Height {
		if err = copier.Copy(&val2, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}

		if err = db.Table("unbonding_delegation").
			Where("validator_address = ? AND delegator_address = ?",
				val2.ValidatorAddress,
				val2.DelegatorAddress,
			).Updates(&updates).Error; err != nil {

			return errors.Wrap(err, "update unbonding_delegation error")
		}

	}

	// v.db.GetGormDB(ctx).Table("unbonding_delegation").Create(val2)

	return nil
}
