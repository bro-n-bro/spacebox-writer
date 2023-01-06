package bank

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func AccountBalanceHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.AccountBalance{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	coinsBytes, err := jsoniter.Marshal(val.Coins)
	if err != nil {
		return err
	}

	val2 := storageModel.AccountBalance{
		Coins:   string(coinsBytes),
		Address: val.Address,
		Height:  val.Height,
	}

	if err = ch.GetGormDB(ctx).Table("account_balance").Create(val2).Error; err != nil {
		return err
	}

	return nil
}
