package core

import (
	"context"
	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
)

func BlockHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.Block{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	val2 := storageModel.Block{
		Height:          val.Height,
		Hash:            val.Hash,
		NumTXS:          int64(val.NumTxs),
		TotalGas:        int64(val.TotalGas),
		ProposerAddress: val.ProposerAddress,
		Timestamp:       val.Timestamp,
	}

	if err := ch.GetGormDB(ctx).Table("block").Create(val2).Error; err != nil {
		return err
	}

	return nil
}
