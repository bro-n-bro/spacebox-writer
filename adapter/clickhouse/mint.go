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
func (ch *Clickhouse) MintParams(vals []model.MintParams) (err error) {
	var (
		params string
	)

	batch := make([]storageModel.MintParams, len(vals))
	for i, val := range vals {
		if params, err = jsoniter.MarshalToString(val.Params); err != nil {
			return err
		}
		batch[i] = storageModel.MintParams{
			Params: params,
			Height: val.Height,
		}
	}

	return ch.gorm.Table(tableMintParams).CreateInBatches(batch, len(batch)).Error
}
