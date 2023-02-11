package clickhouse

import (
	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableMintParams = "mint_params"
)

// MintParams is a method for saving mint params data to clickhouse
func (ch *Clickhouse) MintParams(val model.MintParams) (err error) {
	var (
		paramsBytes []byte
	)

	if paramsBytes, err = jsoniter.Marshal(val.Params); err != nil {
		return err
	}
	return ch.gorm.Table(tableMintParams).Create(storageModel.MintParams{
		Params: string(paramsBytes),
		Height: val.Height,
	}).Error
}
