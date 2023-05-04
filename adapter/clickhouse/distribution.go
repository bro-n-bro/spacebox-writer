package clickhouse

import (
	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableCommunityPool                      = "community_pool"
	tableDistributionParams                 = "distribution_params"
	tableDistributionCommission             = "distribution_commission"
	tableDelegationRewardMessage            = "delegation_reward_message"
	tableProposerReward                     = "proposer_reward"
	tableWithdrawValidatorCommissionMessage = "withdraw_validator_commission_message"
)

// ProposerReward is a method for saving proposer reward data to clickhouse
func (ch *Clickhouse) ProposerReward(vals []model.ProposerReward) (err error) {
	var (
		coins string
	)

	batch := make([]storageModel.ProposerReward, len(vals))

	for i, val := range vals {
		if coins, err = jsoniter.MarshalToString(val.Reward); err != nil {
			return err
		}
		batch[i] = storageModel.ProposerReward{
			Reward:    coins,
			Height:    val.Height,
			Validator: val.Validator,
		}
	}

	return ch.gorm.Table(tableProposerReward).CreateInBatches(batch, len(batch)).Error
}

// CommunityPool is a method for saving community pool data to clickhouse
func (ch *Clickhouse) CommunityPool(vals []model.CommunityPool) (err error) {
	var (
		coins string
	)

	batch := make([]storageModel.CommunityPool, len(vals))

	for i, val := range vals {
		if coins, err = jsoniter.MarshalToString(val.Coins); err != nil {
			return err
		}
		batch[i] = storageModel.CommunityPool{
			Coins:  coins,
			Height: val.Height,
		}
	}

	return ch.gorm.Table(tableCommunityPool).CreateInBatches(batch, len(batch)).Error
}

// DelegationRewardMessage is a method for saving delegation reward message data to clickhouse
func (ch *Clickhouse) DelegationRewardMessage(vals []model.DelegationRewardMessage) (err error) {
	var (
		coins string
	)

	batch := make([]storageModel.DelegationRewardMessage, len(vals))
	for i, val := range vals {
		if coins, err = jsoniter.MarshalToString(val.Coins); err != nil {
			return err
		}
		batch[i] = storageModel.DelegationRewardMessage{
			Coins:            coins,
			Height:           val.Height,
			DelegatorAddress: val.DelegatorAddress,
			ValidatorAddress: val.ValidatorAddress,
			TxHash:           val.TxHash,
			MsgIndex:         val.MsgIndex,
		}
	}

	return ch.gorm.Table(tableDelegationRewardMessage).CreateInBatches(batch, len(batch)).Error
}

// DistributionParams is a method for saving distribution params data to clickhouse
func (ch *Clickhouse) DistributionParams(vals []model.DistributionParams) (err error) {
	var (
		params string
	)

	batch := make([]storageModel.DistributionParams, len(vals))
	for i, val := range vals {
		if params, err = jsoniter.MarshalToString(val.Params); err != nil {
			return err
		}
		batch[i] = storageModel.DistributionParams{
			Params: params,
			Height: val.Height,
		}
	}

	return ch.gorm.Table(tableDistributionParams).CreateInBatches(batch, len(batch)).Error
}

// DistributionCommission is a method for saving distribution commission data to clickhouse
func (ch *Clickhouse) DistributionCommission(vals []model.DistributionCommission) (err error) {
	var (
		amount string
	)

	batch := make([]storageModel.DistributionCommission, len(vals))
	for i, val := range vals {
		if amount, err = jsoniter.MarshalToString(val.Amount); err != nil {
			return err
		}

		batch[i] = storageModel.DistributionCommission{
			Validator: val.Validator,
			Amount:    amount,
			Height:    val.Height,
		}
	}

	return ch.gorm.Table(tableDistributionCommission).CreateInBatches(batch, len(batch)).Error
}

// WithdrawValidatorCommissionMessage is a method for saving withdraw validator commission message data to clickhouse
func (ch *Clickhouse) WithdrawValidatorCommissionMessage(vals []model.WithdrawValidatorCommissionMessage) (err error) {
	var (
		commission string
	)

	batch := make([]storageModel.WithdrawValidatorCommissionMessage, len(vals))
	for i, val := range vals {
		if commission, err = jsoniter.MarshalToString(val.WithdrawCommission); err != nil {
			return err
		}

		batch[i] = storageModel.WithdrawValidatorCommissionMessage{
			ValidatorAddress:   val.ValidatorAddress,
			TxHash:             val.TxHash,
			WithdrawCommission: commission,
			MsgIndex:           val.MsgIndex,
			Height:             val.Height,
		}
	}

	return ch.gorm.Table(tableWithdrawValidatorCommissionMessage).CreateInBatches(batch, len(batch)).Error
}
