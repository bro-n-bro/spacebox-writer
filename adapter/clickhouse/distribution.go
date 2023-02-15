package clickhouse

import (
	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableCommunityPool           = "community_pool"
	tableDistributionParams      = "distribution_params"
	tableDelegationRewardMessage = "delegation_reward_message"
)

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
		params string
	)

	batch := make([]storageModel.DelegationRewardMessage, len(vals))
	for i, val := range vals {
		if params, err = jsoniter.MarshalToString(val.Coins); err != nil {
			return err
		}
		batch[i] = storageModel.DelegationRewardMessage{
			Coins:            params,
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
