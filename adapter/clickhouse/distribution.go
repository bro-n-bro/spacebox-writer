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
func (ch *Clickhouse) CommunityPool(val model.CommunityPool) (err error) {
	var (
		coinsBytes []byte
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}
	return ch.gorm.Table(tableCommunityPool).Create(storageModel.CommunityPool{
		Coins:  string(coinsBytes),
		Height: val.Height,
	}).Error
}

// DelegationRewardMessage is a method for saving delegation reward message data to clickhouse
func (ch *Clickhouse) DelegationRewardMessage(val model.DelegationRewardMessage) (err error) {
	var (
		coinBytes []byte
	)

	if coinBytes, err = jsoniter.Marshal(val.Coin); err != nil {
		return err
	}

	return ch.gorm.Table(tableDelegationRewardMessage).Create(storageModel.DelegationRewardMessage{
		Coin:             string(coinBytes),
		Height:           val.Height,
		DelegatorAddress: val.DelegatorAddress,
		ValidatorAddress: val.ValidatorAddress,
		TxHash:           val.TxHash,
		MsgIndex:         val.MsgIndex,
	}).Error
}

// DistributionParams is a method for saving distribution params data to clickhouse
func (ch *Clickhouse) DistributionParams(val model.DistributionParams) (err error) {
	var (
		paramsBytes []byte
	)

	if paramsBytes, err = jsoniter.Marshal(val.Params); err != nil {
		return err
	}

	return ch.gorm.Table(tableDistributionParams).Create(storageModel.DistributionParams{
		Params: string(paramsBytes),
		Height: val.Height,
	}).Error
}
