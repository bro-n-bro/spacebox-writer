package clickhouse

import (
	"database/sql"

	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableAccountBalance = "account_balance"
	tableSendMessage    = "send_message"
	tableSupply         = "supply"

	insertMultiSendMessage = `
		INSERT INTO spacebox.multisend_message (height, address_from, addresses_to, tx_hash, coins, msg_index)`
)

// AccountBalance is a method for saving account balance data to clickhouse
func (ch *Clickhouse) AccountBalance(vals []model.AccountBalance) (err error) {
	var (
		coins string
	)

	batch := make([]storageModel.AccountBalance, len(vals))
	for i, val := range vals {
		if coins, err = jsoniter.MarshalToString(val.Coins); err != nil {
			return err
		}

		batch[i] = storageModel.AccountBalance{
			Coins:   coins,
			Address: val.Address,
			Height:  val.Height,
		}
	}

	return ch.gorm.Table(tableAccountBalance).CreateInBatches(batch, len(batch)).Error
}

// MultiSendMessage is a method for saving multisend message data to clickhouse
func (ch *Clickhouse) MultiSendMessage(vals []model.MultiSendMessage) (err error) {
	var (
		stmt  *sql.Stmt
		tx    *sql.Tx
		coins string
	)

	if tx, err = ch.sql.Begin(); err != nil {
		return err
	}

	if stmt, err = tx.Prepare(insertMultiSendMessage); err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()

	for _, val := range vals {
		if coins, err = jsoniter.MarshalToString(val.Coins); err != nil {
			return err
		}
		if _, err = stmt.Exec(
			val.Height,
			val.AddressFrom,
			val.AddressesTo,
			val.TxHash,
			coins,
			val.MsgIndex,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// SendMessage is a method for saving send message data to clickhouse
func (ch *Clickhouse) SendMessage(vals []model.SendMessage) (err error) {
	var (
		coins string
	)

	batch := make([]storageModel.SendMessage, len(vals))
	for i, val := range vals {
		if coins, err = jsoniter.MarshalToString(val.Coins); err != nil {
			return err
		}
		batch[i] = storageModel.SendMessage{
			Coins:       coins,
			AddressFrom: val.AddressFrom,
			AddressTo:   val.AddressTo,
			TxHash:      val.TxHash,
			Height:      val.Height,
			MsgIndex:    val.MsgIndex,
		}
	}

	return ch.gorm.Table(tableSendMessage).CreateInBatches(batch, len(batch)).Error
}

// Supply is a method for saving supply data to clickhouse
func (ch *Clickhouse) Supply(vals []model.Supply) (err error) {
	var (
		coins string
	)

	batch := make([]storageModel.Supply, len(vals))
	for i, val := range vals {
		if coins, err = jsoniter.MarshalToString(val.Coins); err != nil {
			return err
		}
		batch[i] = storageModel.Supply{
			Coins:  coins,
			Height: val.Height,
		}
	}

	return ch.gorm.Table(tableSupply).CreateInBatches(batch, len(batch)).Error
}
