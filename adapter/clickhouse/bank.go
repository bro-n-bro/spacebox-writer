package clickhouse

import (
	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/hexy-dev/spacebox-writer/adapter/clickhouse/models"
	"github.com/hexy-dev/spacebox/broker/model"
)

func (ch *Clickhouse) AccountBalance(val model.AccountBalance) (err error) {
	var (
		coinsBytes []byte
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	if err = ch.gorm.Table("account_balance").Create(storageModel.AccountBalance{
		Coins:   string(coinsBytes),
		Address: val.Address,
		Height:  val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) MultisendMessage(val model.MultiSendMessage) (err error) {
	var (
		coinsBytes       []byte
		addressesToBytes []byte
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	if addressesToBytes, err = jsoniter.Marshal(val.AddressesTo); err != nil {
		return err
	}

	if err = ch.gorm.Table("multisend_message").Create(storageModel.MultiSendMessage{
		Coins:       string(coinsBytes),
		AddressesTo: string(addressesToBytes),
		AddressFrom: val.AddressFrom,
		TxHash:      val.TxHash,
		Height:      val.Height,
		MsgIndex:    val.MsgIndex,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) SendMessage(val model.SendMessage) (err error) {
	var (
		coinsBytes []byte
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	if err = ch.gorm.Table("send_message").Create(storageModel.SendMessage{
		Coins:       string(coinsBytes),
		AddressFrom: val.AddressFrom,
		AddressTo:   val.AddressTo,
		TxHash:      val.TxHash,
		Height:      val.Height,
		MsgIndex:    val.MsgIndex,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) Supply(val model.Supply) (err error) {
	var (
		coinsBytes []byte
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	if err = ch.gorm.Table("supply").Create(storageModel.Supply{
		Coins:  string(coinsBytes),
		Height: val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}
