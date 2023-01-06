package bank

import (
	"context"
	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
)

func SupplyHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.Supply{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	coinsBytes, err := jsoniter.Marshal(val.Coins)
	if err != nil {
		return err
	}

	val2 := storageModel.Supply{
		Coins:  string(coinsBytes),
		Height: val.Height,
	}

	if err = ch.GetGormDB(ctx).Table("supply").Create(val2).Error; err != nil {
		return err
	}

	return nil
}