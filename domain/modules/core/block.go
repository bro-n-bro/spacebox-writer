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

	val := model.AccountBalance{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	coinsBytes, err := jsoniter.Marshal(val.Coins)
	if err != nil {
		return err
	}

	val2 := storageModel.AccountBalance{
		Address: val.Address,
		Coins:   string(coinsBytes),
		Height:  val.Height,
	}

	ch.GetGormDB(ctx).Table("account_balance").Create(val2)

	return nil
}
