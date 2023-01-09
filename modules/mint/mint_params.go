package mint

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func MintParamsHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.MintParams{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	paramsBytes, err := jsoniter.Marshal(val.Params)
	if err != nil {
		return err
	}

	if err = ch.GetGormDB(ctx).Table("mint_params").Create(storageModel.MintParams{
		Height: val.Height,
		Params: string(paramsBytes),
	}).Error; err != nil {
		return err
	}

	return nil
}
