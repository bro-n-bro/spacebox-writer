package clickhouse

import (
	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableCommunityPool           = "community_pool"
	tableDistributionParams      = "distribution_params"
	tableDelegationRewardMessage = "delegation_reward_message"
	tableValidatorCommission     = "validator_commission"
)

func (ch *Clickhouse) CommunityPool(val model.CommunityPool) (err error) {
	var (
		coinsBytes     []byte
		updates        storageModel.CommunityPool
		prevValStorage storageModel.CommunityPool
		valStorage     storageModel.CommunityPool
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	valStorage = storageModel.CommunityPool{
		Coins:  string(coinsBytes),
		Height: val.Height,
	}

	if err = ch.gorm.Table(tableCommunityPool).
		First(&prevValStorage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table(tableCommunityPool).Create(val).Error; err != nil {
				return errors.Wrap(err, "error of create")
			}
			return nil
		}
		return errors.Wrap(err, "error of database: %w")
	}

	if valStorage.Height > prevValStorage.Height {
		if err = copier.Copy(&valStorage, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}
		if err = ch.gorm.Table(tableCommunityPool).
			Where("height = ?", valStorage.Height).
			Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "error of update")
		}
	}

	return nil
}

func (ch *Clickhouse) DelegationRewardMessage(val model.DelegationRewardMessage) (err error) {
	var (
		paramsBytes []byte
		exists      bool
	)

	if paramsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	if exists, err = ch.ExistsTx(tableDelegationRewardMessage, val.TxHash, val.MsgIndex); err != nil {
		return err
	}

	if !exists {
		if err = ch.gorm.Table(tableDelegationRewardMessage).Create(storageModel.DelegationRewardMessage{
			Coins:            string(paramsBytes),
			Height:           val.Height,
			DelegatorAddress: val.DelegatorAddress,
			ValidatorAddress: val.ValidatorAddress,
			TxHash:           val.TxHash,
			MsgIndex:         val.MsgIndex,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

func (ch *Clickhouse) DistributionParams(val model.DistributionParams) (err error) {
	var (
		paramsBytes    []byte
		updates        storageModel.DistributionParams
		prevValStorage storageModel.DistributionParams
		valStorage     storageModel.DistributionParams
	)

	if paramsBytes, err = jsoniter.Marshal(val.Params); err != nil {
		return err
	}

	valStorage = storageModel.DistributionParams{
		Params: string(paramsBytes),
		Height: val.Height,
	}

	if err = ch.gorm.Table(tableDistributionParams).
		First(&prevValStorage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table(tableDistributionParams).Create(val).Error; err != nil {
				return errors.Wrap(err, "error of create")
			}
			return nil
		}
		return errors.Wrap(err, "error of database: %w")
	}

	if valStorage.Height > prevValStorage.Height {
		if err = copier.Copy(&valStorage, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}
		if err = ch.gorm.Table(tableDistributionParams).
			Where("height = ?", valStorage.Height).
			Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "error of update")
		}
	}

	return nil
}

func (ch *Clickhouse) ValidatorCommission(val model.ValidatorCommission) (err error) {
	var (
		updates        model.ValidatorCommission
		prevValStorage model.ValidatorCommission
	)

	if err = ch.gorm.Table(tableValidatorCommission).
		Where("operator_address = ?", val.OperatorAddress).
		First(&prevValStorage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table(tableValidatorCommission).Create(val).Error; err != nil {
				return errors.Wrap(err, "error of create")
			}
			return nil
		}
		return errors.Wrap(err, "error of database: %w")
	}

	if val.Height > prevValStorage.Height {
		if err = copier.Copy(&val, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}
		if err = ch.gorm.Table(tableValidatorCommission).
			Where("operator_address = ?", val.OperatorAddress).
			Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "error of update")
		}
	}

	return nil
}
