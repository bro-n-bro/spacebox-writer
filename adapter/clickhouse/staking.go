package clickhouse

import (
	"database/sql"
	"encoding/json"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableDelegation                 = "delegation"
	tableDelegationMessage          = "delegation_message"
	tableRedelegation               = "redelegation"
	tableRedelegationMessage        = "redelegation_message"
	tableUnbondingDelegation        = "unbonding_delegation"
	tableUnbondingDelegationMessage = "unbonding_delegation_message"
	tableValidator                  = "validator"
	tableValidatorDescription       = "validator_description"
	tableValidatorStatus            = "validator_status"
	tableValidatorInfo              = "validator_info"
	tableStackingPool               = "stacking_pool"
	tableStackingParams             = "stacking_params"
)

func (ch *Clickhouse) Delegation(val model.Delegation) (err error) {
	var (
		coinBytes      []byte
		updates        storageModel.Delegation
		prevValStorage storageModel.Delegation
		valStorage     storageModel.Delegation
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return errors.Wrap(err, "marshall error")
	}

	valStorage = storageModel.Delegation{
		OperatorAddress:  val.OperatorAddress,
		DelegatorAddress: val.DelegatorAddress,
		Coin:             string(coinBytes),
		Height:           val.Height,
	}

	if err = ch.gorm.Table(tableDelegation).
		Where("operator_address = ? AND delegator_address = ?", valStorage.OperatorAddress, valStorage.DelegatorAddress).
		First(&prevValStorage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table(tableDelegation).Create(valStorage).Error; err != nil {
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
		if err = ch.gorm.Table(tableDelegation).
			Where("operator_address = ? AND delegator_address = ?", valStorage.OperatorAddress, valStorage.DelegatorAddress).
			Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "error of update")
		}
	}

	return nil
}

func (ch *Clickhouse) DelegationMessage(val model.DelegationMessage) (err error) {
	var (
		coinBytes []byte
		exists    bool
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return errors.Wrap(err, "marshall error")
	}

	if exists, err = ch.ExistsTx(tableDelegationMessage, val.TxHash, val.MsgIndex); err != nil {
		return err
	}

	if !exists {
		if err = ch.gorm.Table(tableDelegationMessage).Create(storageModel.DelegationMessage{
			OperatorAddress:  val.OperatorAddress,
			DelegatorAddress: val.DelegatorAddress,
			Coin:             string(coinBytes),
			Height:           val.Height,
			TxHash:           val.TxHash,
			MsgIndex:         val.MsgIndex,
		}).Error; err != nil {
			return err
		}
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

	if err = ch.gorm.Table(tableRedelegation).Create(storageModel.Redelegation{
		// nolint: lll
		// DEFAULT CLICKHOUSE ZERO TIME: 1970-01-01 00:00:00.
		// https://github.com/ClickHouse/ClickHouse/issues/15081
		// Source driver code: https://github.com/ClickHouse/clickhouse-go/blob/0c79b0fd50ee848198dde5a3392850820974ee6f/lib/column/datetime.go#L200
		// Driver issue: https://github.com/ClickHouse/clickhouse-go/issues/882
		CompletionTime: sql.NullTime{
			Time:  val.CompletionTime,
			Valid: !val.CompletionTime.IsZero(),
		},
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
		exists    bool
	)

	if coinBytes, err = json.Marshal(val.Coin); err != nil {
		return err
	}

	if exists, err = ch.ExistsTx(tableRedelegationMessage, val.TxHash, val.MsgIndex); err != nil {
		return err
	}

	if !exists {
		if err = ch.gorm.Table(tableRedelegationMessage).Create(storageModel.RedelegationMessage{
			// nolint: lll
			// DEFAULT CLICKHOUSE ZERO TIME: 1970-01-01 00:00:00.
			// https://github.com/ClickHouse/ClickHouse/issues/15081
			// Source driver code: https://github.com/ClickHouse/clickhouse-go/blob/0c79b0fd50ee848198dde5a3392850820974ee6f/lib/column/datetime.go#L200
			// Driver issue: https://github.com/ClickHouse/clickhouse-go/issues/882
			CompletionTime: sql.NullTime{
				Time:  val.CompletionTime,
				Valid: !val.CompletionTime.IsZero(),
			},
			Coin:                string(coinBytes),
			DelegatorAddress:    val.DelegatorAddress,
			SrcValidatorAddress: val.SrcValidatorAddress,
			DstValidatorAddress: val.DstValidatorAddress,
			Height:              val.Height,
			TxHash:              val.TxHash,
		}).Error; err != nil {
			return err
		}

	}

	return nil
}

func (ch *Clickhouse) StakingParams(val model.StakingParams) (err error) {
	var (
		prevStorageVal storageModel.StakingParams
		storageVal     storageModel.StakingParams
		updates        storageModel.StakingParams
		paramsBytes    []byte
	)

	if paramsBytes, err = jsoniter.Marshal(val.Params); err != nil {
		return err
	}

	storageVal = storageModel.StakingParams{
		Params: string(paramsBytes),
		Height: val.Height,
	}

	if err = ch.gorm.Table(tableStackingParams).
		Where("height = ?", val.Height).
		First(&prevStorageVal).Error; !errors.Is(gorm.ErrRecordNotFound, err) {
		return err
	} else if storageVal.Params != prevStorageVal.Params {
		if err = copier.Copy(&storageVal, &updates); err != nil {
			return err
		}
		if err = ch.gorm.Table(tableStackingParams).
			Where("height = ?", storageVal.Height).
			Updates(&updates).Error; err != nil {
			return err
		}
	}

	return nil
}

func (ch *Clickhouse) StakingPool(val model.StakingPool) (err error) {
	return ch.gorm.Table(tableStackingPool).Create(val).Error
}

func (ch *Clickhouse) UnbondingDelegation(val model.UnbondingDelegation) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return err
	}

	if err = ch.gorm.Table(tableUnbondingDelegation).Create(storageModel.UnbondingDelegation{
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
		exists    bool
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return err
	}

	if exists, err = ch.ExistsTx(tableUnbondingDelegationMessage, val.TxHash, val.MsgIndex); err != nil {
		return err
	}

	if !exists {
		if err = ch.gorm.Table(tableUnbondingDelegationMessage).
			Create(storageModel.UnbondingDelegationMessage{
				CompletionTimestamp: val.CompletionTimestamp,
				Coin:                string(coinBytes),
				DelegatorAddress:    val.DelegatorAddress,
				ValidatorAddress:    val.ValidatorAddress,
				Height:              val.Height,
				TxHash:              val.TxHash,
				MsgIndex:            val.MsgIndex,
			}).Error; err != nil {
			return err
		}
	}
	return nil
}

// Validator TODO: добавить height
func (ch *Clickhouse) Validator(val model.Validator) (err error) {
	//var count int66
	//if ch.gorm.Table(tableValidator).
	//	Where("consensus_address = ?", val.ConsensusAddress).
	//	Count(&count); count != 0 {
	//	return nil
	//}

	if err = ch.gorm.Table(tableValidator).Create(val).Error; err != nil {
		return errors.Wrap(err, "create validator error")
	}

	return nil
}

func (ch *Clickhouse) ValidatorInfo(val model.ValidatorInfo) (err error) {
	var (
		updates  model.ValidatorInfo
		previous model.ValidatorInfo
	)

	if err = ch.gorm.Table(tableValidatorInfo).
		Where("consensus_address = ? AND operator_address = ?", val.ConsensusAddress, val.OperatorAddress).
		First(&previous).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table(tableValidatorInfo).Create(val).Error; err != nil {
				return errors.Wrap(err, "create validator_info error")
			}
			return nil
		}
		return err
	}

	if val.Height > previous.Height {
		if err = copier.Copy(&val, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}
		if err = ch.gorm.Table(tableValidatorInfo).
			Where("consensus_address = ? AND operator_address = ?", val.ConsensusAddress, val.OperatorAddress).
			Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "update validator_info error")
		}
	}

	return nil
}

func (ch *Clickhouse) ValidatorStatus(val model.ValidatorStatus) (err error) {
	var (
		updates  model.ValidatorStatus
		previous model.ValidatorStatus
	)

	if err = ch.gorm.Table(tableValidatorStatus).
		Where("validator_address = ?", val.ValidatorAddress).
		First(&previous).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table(tableValidatorStatus).Create(val).Error; err != nil {
				return errors.Wrap(err, "create validator_status error")
			}
			return nil
		}
		return err
	}

	if val.Height > previous.Height {
		if err = copier.Copy(&val, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}
		if err = ch.gorm.Table(tableValidatorStatus).
			Where("validator_address = ?", val.ValidatorAddress).
			Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "update validator_status error")
		}
	}

	return nil
}

func (ch *Clickhouse) ValidatorDescription(val model.ValidatorDescription) (err error) {
	var (
		updates  model.ValidatorDescription
		previous model.ValidatorDescription
	)

	if err = ch.gorm.Table(tableValidatorDescription).
		Where("operator_address = ?", val.OperatorAddress).
		First(&previous).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table(tableValidatorDescription).Create(val).Error; err != nil {
				return errors.Wrap(err, "create validator_description error")
			}
			return nil
		}
		return err
	}

	if val.Height > previous.Height {
		if err = copier.Copy(&val, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}
		if err = ch.gorm.Table(tableValidatorDescription).
			Where("operator_address = ?", val.OperatorAddress).
			Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "update validator_description error")
		}
	}

	return nil
}
