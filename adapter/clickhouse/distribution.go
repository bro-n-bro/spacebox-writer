package clickhouse

import (
	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/hexy-dev/spacebox-writer/adapter/clickhouse/models"
	"github.com/hexy-dev/spacebox/broker/model"
)

func (ch *Clickhouse) CommunityPool(val model.CommunityPool) (err error) {
	var (
		coinsBytes []byte
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	if err = ch.gorm.Table("community_pool").Create(storageModel.CommunityPool{
		Coins:  string(coinsBytes),
		Height: val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) DelegationRewardMessage(val model.DelegationRewardMessage) (err error) {
	var (
		paramsBytes []byte
	)

	if paramsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	if err = ch.gorm.Table("delegation_reward_message").Create(storageModel.DelegationRewardMessage{
		Coins:            string(paramsBytes),
		Height:           val.Height,
		DelegatorAddress: val.DelegatorAddress,
		ValidatorAddress: val.ValidatorAddress,
		TxHash:           val.TxHash,
		MsgIndex:         val.MsgIndex,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) DistributionParams(val model.DistributionParams) (err error) {
	var (
		paramsBytes []byte
	)

	if paramsBytes, err = jsoniter.Marshal(val.Params); err != nil {
		return err
	}

	if err = ch.gorm.Table("distribution_params").Create(storageModel.DistributionParams{
		Params: string(paramsBytes),
		Height: val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) ValidatorCommission(val model.ValidatorCommission) (err error) {
	if err = ch.gorm.Table("validator_commission").Create(val).Error; err != nil {
		return err
	}

	return nil
}
