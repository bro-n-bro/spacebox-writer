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
		INSERT INTO spacebox.multisend_message (height, address_from, addresses_to, tx_hash, coins, msg_index)
		VALUES (?, ?, ?, ?, ?, ?);`
)

// AccountBalance is a method for saving account balance data to clickhouse
func (ch *Clickhouse) AccountBalance(val model.AccountBalance) (err error) {
	var (
		coinsBytes []byte
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}
	if err = ch.gorm.Table(tableAccountBalance).Create(storageModel.AccountBalance{
		Coins:   string(coinsBytes),
		Address: val.Address,
		Height:  val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}

// MultiSendMessage is a method for saving multisend message data to clickhouse
func (ch *Clickhouse) MultiSendMessage(val model.MultiSendMessage) (err error) {
	var (
		stmt       *sql.Stmt
		tx         *sql.Tx
		coinsBytes []byte
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	if tx, err = ch.sql.Begin(); err != nil {
		return err
	}
	if stmt, err = tx.Prepare(insertMultiSendMessage); err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()
	if _, err = stmt.Exec(
		val.Height,
		val.AddressFrom,
		val.AddressesTo,
		val.TxHash,
		string(coinsBytes),
		val.MsgIndex,
	); err != nil {
		return err
	}
	return tx.Commit()
}

// SendMessage is a method for saving send message data to clickhouse
func (ch *Clickhouse) SendMessage(val model.SendMessage) (err error) {
	var (
		coinsBytes []byte
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}
	if err = ch.gorm.Table(tableSendMessage).
		Create(storageModel.SendMessage{
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

// Supply is a method for saving supply data to clickhouse
func (ch *Clickhouse) Supply(val model.Supply) (err error) {
	var (
		coinsBytes []byte
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}
	if err = ch.gorm.Table(tableSupply).Create(storageModel.Supply{
		Coins:  string(coinsBytes),
		Height: val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}
