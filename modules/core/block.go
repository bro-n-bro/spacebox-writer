package core

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func BlockHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.Block{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	if err := ch.GetGormDB(ctx).Table("block").Create(storageModel.Block{
		Height:          val.Height,
		Hash:            val.Hash,
		NumTXS:          int64(val.NumTxs),
		TotalGas:        int64(val.TotalGas),
		ProposerAddress: val.ProposerAddress,
		Timestamp:       val.Timestamp,
	}).Error; err != nil {
		return err
	}

	return nil
}
