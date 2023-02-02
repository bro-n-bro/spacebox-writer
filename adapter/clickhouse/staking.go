package clickhouse

import (
	"database/sql"
	"encoding/json"

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
)

func (ch *Clickhouse) Delegation(val model.Delegation) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return errors.Wrap(err, "marshall error")
	}

	return ch.gorm.Table(tableDelegation).Create(storageModel.Delegation{
		OperatorAddress:  val.OperatorAddress,
		DelegatorAddress: val.DelegatorAddress,
		Coin:             string(coinBytes),
		Height:           val.Height,
	}).Error
}

func (ch *Clickhouse) DelegationMessage(val model.DelegationMessage) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return errors.Wrap(err, "marshall error")
	}

	return ch.gorm.Table(tableDelegationMessage).Create(storageModel.DelegationMessage{
		OperatorAddress:  val.OperatorAddress,
		DelegatorAddress: val.DelegatorAddress,
		Coin:             string(coinBytes),
		Height:           val.Height,
		TxHash:           val.TxHash,
		MsgIndex:         val.MsgIndex,
	}).Error
}

func (ch *Clickhouse) Redelegation(val model.Redelegation) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return errors.Wrap(err, "marshall error")
	}

	return ch.gorm.Table(tableRedelegation).Create(storageModel.Redelegation{
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
	}).Error
}

func (ch *Clickhouse) RedelegationMessage(val model.RedelegationMessage) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = json.Marshal(val.Coin); err != nil {
		return err
	}

	return ch.gorm.Table(tableRedelegationMessage).Create(storageModel.RedelegationMessage{
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
	}).Error
}

func (ch *Clickhouse) StakingParams(val model.StakingParams) (err error) {
	var (
		paramsBytes []byte
	)

	if paramsBytes, err = jsoniter.Marshal(val.Params); err != nil {
		return err
	}

	return ch.gorm.Table(tableStakingParams).Create(storageModel.StakingParams{
		Params: string(paramsBytes),
		Height: val.Height,
	}).Error
}

func (ch *Clickhouse) UnbondingDelegation(val model.UnbondingDelegation) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return err
	}

	return ch.gorm.Table(tableUnbondingDelegation).Create(storageModel.UnbondingDelegation{
		CompletionTimestamp: val.CompletionTimestamp,
		Coin:                string(coinBytes),
		DelegatorAddress:    val.DelegatorAddress,
		ValidatorAddress:    val.ValidatorAddress,
		Height:              val.Height,
	}).Error
}

func (ch *Clickhouse) UnbondingDelegationMessage(val model.UnbondingDelegationMessage) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return err
	}

	return ch.gorm.Table(tableUnbondingDelegationMessage).
		Create(storageModel.UnbondingDelegationMessage{
			CompletionTimestamp: val.CompletionTimestamp,
			Coin:                string(coinBytes),
			DelegatorAddress:    val.DelegatorAddress,
			ValidatorAddress:    val.ValidatorAddress,
			Height:              val.Height,
			TxHash:              val.TxHash,
			MsgIndex:            val.MsgIndex,
		}).Error
}
