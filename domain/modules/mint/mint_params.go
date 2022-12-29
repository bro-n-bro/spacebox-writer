package mint

import (
	"context"
	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
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

	val2 := storageModel.MintParams{
		Height: val.Height,
		Params: string(paramsBytes),
	}

	ch.GetGormDB(ctx).Table("mint_params").Create(val2)

	return nil
}
