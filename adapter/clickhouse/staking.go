package clickhouse

import (
	"encoding/json"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"
)

func (ch *Clickhouse) Delegation(val model.Delegation) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return errors.Wrap(err, "marshall error")
	}

	if err = ch.gorm.Table("delegation").Create(storageModel.Delegation{
		OperatorAddress:  val.OperatorAddress,
		DelegatorAddress: val.DelegatorAddress,
		Coin:             string(coinBytes),
		Height:           val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) DelegationMessage(val model.DelegationMessage) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return errors.Wrap(err, "marshall error")
	}

	if err = ch.gorm.Table("delegation_message").Create(storageModel.DelegationMessage{
		OperatorAddress:  val.OperatorAddress,
		DelegatorAddress: val.DelegatorAddress,
		Coin:             string(coinBytes),
		Height:           val.Height,
		TxHash:           val.TxHash,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) Redelegation(val model.Redelegation) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return errors.Wrap(err, "marshall error")
	}

	if err = ch.gorm.Table("redelegation").Create(storageModel.Redelegation{
		CompletionTime:      val.CompletionTime,
		Coin:                string(coinBytes),
		DelegatorAddress:    val.DelegatorAddress,
		SrcValidatorAddress: val.SrcValidatorAddress,
		DstValidatorAddress: val.DstValidatorAddress,
		Height:              val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) RedelegationMessage(val model.RedelegationMessage) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = json.Marshal(val.Coin); err != nil {
		return err
	}

	if err = ch.gorm.Table("redelegation_message").Create(storageModel.RedelegationMessage{
		CompletionTime:      val.CompletionTime,
		Coin:                string(coinBytes),
		DelegatorAddress:    val.DelegatorAddress,
		SrcValidatorAddress: val.SrcValidatorAddress,
		DstValidatorAddress: val.DstValidatorAddress,
		Height:              val.Height,
		TxHash:              val.TxHash,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) StakingParams(val model.StakingParams) (err error) {
	var (
		stVal       storageModel.StakingParams
		updates     storageModel.StakingParams
		paramsBytes []byte
	)

	if paramsBytes, err = jsoniter.Marshal(val.Params); err != nil {
		return err
	}

	val2 := storageModel.StakingParams{
		Params: string(paramsBytes),
		Height: val.Height,
	}

	err = ch.gorm.Table("staking_params").Where("height = ?", val.Height).First(&stVal).Error
	if !errors.Is(gorm.ErrRecordNotFound, err) {
		return err
	} else if val2.Params != stVal.Params {
		if err = copier.Copy(&val2, &updates); err != nil {
			return err
		}
		if err = ch.gorm.Table("staking_params").Where("height = ?", val2.Height).Updates(&updates).
			Error; err != nil {
			return err
		}
	}

	return nil
}

func (ch *Clickhouse) StakingPool(val model.StakingPool) (err error) {
	if err = ch.gorm.Table("staking_pool").Create(val).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) UnbondingDelegation(val model.UnbondingDelegation) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return err
	}

	if err = ch.gorm.Table("unbonding_delegation").Create(storageModel.UnbondingDelegation{
		CompletionTimestamp: val.CompletionTimestamp,
		Coin:                string(coinBytes),
		DelegatorAddress:    val.DelegatorAddress,
		ValidatorAddress:    val.ValidatorAddress,
		Height:              val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) UnbondingDelegationMessage(val model.UnbondingDelegationMessage) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return err
	}

	if err = ch.gorm.Table("unbonding_delegation_message").
		Create(storageModel.UnbondingDelegationMessage{
			CompletionTimestamp: val.CompletionTimestamp,
			Coin:                string(coinBytes),
			DelegatorAddress:    val.DelegatorAddress,
			ValidatorAddress:    val.ValidatorAddress,
			Height:              val.Height,
			TxHash:              val.TxHash,
		}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) Validator(val model.Validator) (err error) {
	var (
		count int64
	)

	if ch.gorm.Table("validator").
		Where("consensus_address = ?", val.ConsensusAddress).
		Count(&count); count != 0 {
		return nil
	}

	if err = ch.gorm.Table("validator").Create(val).Error; err != nil {
		return errors.Wrap(err, "create validator error")
	}

	return nil
}

func (ch *Clickhouse) ValidatorInfo(val model.ValidatorInfo) (err error) {
	var (
		updates model.ValidatorInfo
		getVal  model.ValidatorInfo
	)

	if err = ch.gorm.Table("validator_info").Where("consensus_address = ?", val.ConsensusAddress).
		First(&getVal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table("validator_info").Create(val).Error; err != nil {
				return errors.Wrap(err, "create validator_info error")
			}
			return nil
		}
		return err
	}

	if val.Height < getVal.Height {
		if err = copier.Copy(&val, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}
		if err = ch.gorm.Table("validator_info").
			Where("consensus_address = ?", val.ConsensusAddress).
			Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "update validator_info error")
		}
	}

	return nil
}

func (ch *Clickhouse) ValidatorStatus(val model.ValidatorStatus) (err error) {
	var (
		updates model.ValidatorStatus
		getVal  model.ValidatorStatus
	)

	if err = ch.gorm.Table("validator_status").
		Where("validator_address = ?", val.ValidatorAddress).
		First(&getVal).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table("validator_status").Create(val).Error; err != nil {
				return errors.Wrap(err, "create validator_status erorr")
			}
			return nil
		}
		return err
	}

	if val.Height > getVal.Height {
		if err = copier.Copy(&val, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}
		if err = ch.gorm.Table("validator_status").
			Where("validator_address = ?", val.ValidatorAddress).
			Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "update validator_status error")
		}
	}

	return nil
}

// func (ch *Clickhouse) ValidatorDescription(val model.ValidatorDescription) (err error) {
//	return nil
// }
