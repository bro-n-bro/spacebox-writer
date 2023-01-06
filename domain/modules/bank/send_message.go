package bank

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func SendMessageHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.SendMessage{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	coinsBytes, err := jsoniter.Marshal(val.Coins)
	if err != nil {
		return err
	}

	val2 := storageModel.SendMessage{
		Coins:       string(coinsBytes),
		AddressFrom: val.AddressFrom,
		AddressTo:   val.AddressTo,
		TxHash:      val.TxHash,
		Height:      val.Height,
	}

	if err = ch.GetGormDB(ctx).Table("send_message").Create(val2).Error; err != nil {
		return err
	}

	return nil
}
