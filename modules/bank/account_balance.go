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

	if err = ch.GetGormDB(ctx).Table("account_balance").Create(storageModel.AccountBalance{
		Coins:   string(coinsBytes),
		Address: val.Address,
		Height:  val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}
