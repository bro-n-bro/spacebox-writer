package clickhouse

import (
	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
)

func (ch *Clickhouse) AnnualProvision(val model.AnnualProvision) (err error) {
	if err = ch.gorm.Table("annual_provision").Create(val).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) MintParams(val model.MintParams) (err error) {
	var (
		paramsBytes []byte
	)

	if paramsBytes, err = jsoniter.Marshal(val.Params); err != nil {
		return err
	}

	if err = ch.gorm.Table("mint_params").Create(storageModel.MintParams{
		Height: val.Height,
		Params: string(paramsBytes),
	}).Error; err != nil {
		return err
	}

	return nil
}
