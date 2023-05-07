package clickhouse

import (
	"database/sql"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

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
	tableStakingParams              = "staking_params"
	tableEditValidatorMessage       = "edit_validator_message"
)

// Delegation is a method for saving delegation data to clickhouse
func (ch *Clickhouse) Delegation(vals []model.Delegation) (err error) {
	var (
		coin string
	)
	batch := make([]storageModel.Delegation, len(vals))
	for i, val := range vals {
		if coin, err = jsoniter.MarshalToString(val.Coin); err != nil {
			return errors.Wrap(err, "marshall error")
		}
		batch[i] = storageModel.Delegation{
			OperatorAddress:  val.OperatorAddress,
			DelegatorAddress: val.DelegatorAddress,
			Coin:             coin,
			Height:           val.Height,
		}
	}

	return ch.gorm.Table(tableDelegation).CreateInBatches(batch, len(batch)).Error
}

// DelegationMessage is a method for saving delegation message data to clickhouse
func (ch *Clickhouse) DelegationMessage(vals []model.DelegationMessage) (err error) {
	var (
		coin string
	)

	batch := make([]storageModel.DelegationMessage, len(vals))
	for i, val := range vals {
		if coin, err = jsoniter.MarshalToString(val.Coin); err != nil {
			return errors.Wrap(err, "marshall error")
		}

		batch[i] = storageModel.DelegationMessage{
			OperatorAddress:  val.OperatorAddress,
			DelegatorAddress: val.DelegatorAddress,
			Coin:             coin,
			Height:           val.Height,
			TxHash:           val.TxHash,
			MsgIndex:         val.MsgIndex,
		}
	}

	return ch.gorm.Table(tableDelegationMessage).CreateInBatches(batch, len(batch)).Error
}

// Redelegation is a method for saving redelegation data to clickhouse
func (ch *Clickhouse) Redelegation(vals []model.Redelegation) (err error) {
	var (
		coin string
	)

	batch := make([]storageModel.Redelegation, len(vals))
	for i, val := range vals {
		if coin, err = jsoniter.MarshalToString(val.Coin); err != nil {
			return errors.Wrap(err, "marshall error")
		}
		batch[i] = storageModel.Redelegation{
			// nolint: lll
			// DEFAULT CLICKHOUSE ZERO TIME: 1970-01-01 00:00:00.
			// https://github.com/ClickHouse/ClickHouse/issues/15081
			// Source driver code: https://github.com/ClickHouse/clickhouse-go/blob/0c79b0fd50ee848198dde5a3392850820974ee6f/lib/column/datetime.go#L200
			// Driver issue: https://github.com/ClickHouse/clickhouse-go/issues/882
			CompletionTime: sql.NullTime{
				Time:  val.CompletionTime,
				Valid: !val.CompletionTime.IsZero(),
			},
			Coin:                coin,
			DelegatorAddress:    val.DelegatorAddress,
			SrcValidatorAddress: val.SrcValidatorAddress,
			DstValidatorAddress: val.DstValidatorAddress,
			Height:              val.Height,
		}
	}

	return ch.gorm.Table(tableRedelegation).CreateInBatches(batch, len(batch)).Error
}

// RedelegationMessage is a method for saving redelegation message data to clickhouse
func (ch *Clickhouse) RedelegationMessage(vals []model.RedelegationMessage) (err error) {
	var (
		coin string
	)
	batch := make([]storageModel.RedelegationMessage, len(vals))
	for i, val := range vals {
		if coin, err = jsoniter.MarshalToString(val.Coin); err != nil {
			return err
		}
		batch[i] = storageModel.RedelegationMessage{
			// nolint: lll
			// DEFAULT CLICKHOUSE ZERO TIME: 1970-01-01 00:00:00.
			// https://github.com/ClickHouse/ClickHouse/issues/15081
			// Source driver code: https://github.com/ClickHouse/clickhouse-go/blob/0c79b0fd50ee848198dde5a3392850820974ee6f/lib/column/datetime.go#L200
			// Driver issue: https://github.com/ClickHouse/clickhouse-go/issues/882
			CompletionTime: sql.NullTime{
				Time:  val.CompletionTime,
				Valid: !val.CompletionTime.IsZero(),
			},
			Coin:                coin,
			DelegatorAddress:    val.DelegatorAddress,
			SrcValidatorAddress: val.SrcValidatorAddress,
			DstValidatorAddress: val.DstValidatorAddress,
			Height:              val.Height,
			TxHash:              val.TxHash,
		}
	}

	return ch.gorm.Table(tableRedelegationMessage).CreateInBatches(batch, len(batch)).Error
}

// StakingParams is a method for saving staking params data to clickhouse
func (ch *Clickhouse) StakingParams(vals []model.StakingParams) (err error) {
	var (
		params string
	)

	batch := make([]storageModel.StakingParams, len(vals))
	for i, val := range vals {
		if params, err = jsoniter.MarshalToString(val.Params); err != nil {
			return err
		}
		batch[i] = storageModel.StakingParams{
			Params: params,
			Height: val.Height,
		}
	}

	return ch.gorm.Table(tableStakingParams).CreateInBatches(batch, len(batch)).Error
}

// UnbondingDelegation is a method for saving unbonding delegation data to clickhouse
func (ch *Clickhouse) UnbondingDelegation(vals []model.UnbondingDelegation) (err error) {
	var (
		coin string
	)

	batch := make([]storageModel.UnbondingDelegation, len(vals))
	for i, val := range vals {
		if coin, err = jsoniter.MarshalToString(val.Coin); err != nil {
			return err
		}
		batch[i] = storageModel.UnbondingDelegation{
			CompletionTimestamp: val.CompletionTimestamp,
			Coin:                coin,
			DelegatorAddress:    val.DelegatorAddress,
			ValidatorAddress:    val.ValidatorAddress,
			Height:              val.Height,
		}
	}

	return ch.gorm.Table(tableUnbondingDelegation).CreateInBatches(batch, len(batch)).Error
}

// UnbondingDelegationMessage is a method for saving unbonding delegation message data to clickhouse
func (ch *Clickhouse) UnbondingDelegationMessage(vals []model.UnbondingDelegationMessage) (err error) {
	var (
		coin string
	)

	batch := make([]storageModel.UnbondingDelegationMessage, len(vals))
	for i, val := range vals {
		if coin, err = jsoniter.MarshalToString(val.Coin); err != nil {
			return err
		}

		batch[i] = storageModel.UnbondingDelegationMessage{
			CompletionTimestamp: val.CompletionTimestamp,
			Coin:                coin,
			DelegatorAddress:    val.DelegatorAddress,
			ValidatorAddress:    val.ValidatorAddress,
			Height:              val.Height,
			TxHash:              val.TxHash,
			MsgIndex:            val.MsgIndex,
		}
	}

	return ch.gorm.Table(tableUnbondingDelegationMessage).CreateInBatches(batch, len(batch)).Error
}

// EditValidatorMessage is a method for saving edit validator message data to clickhouse
func (ch *Clickhouse) EditValidatorMessage(vals []model.EditValidatorMessage) (err error) {
	var (
		description string
	)

	batch := make([]storageModel.EditValidatorMessage, len(vals))
	for i, val := range vals {
		if description, err = jsoniter.MarshalToString(val.Description); err != nil {
			return err
		}

		batch[i] = storageModel.EditValidatorMessage{
			Height:      val.Height,
			MsgIndex:    val.Index,
			TxHash:      val.Hash,
			Description: description,
		}
	}

	return ch.gorm.Table(tableEditValidatorMessage).CreateInBatches(batch, len(batch)).Error
}
