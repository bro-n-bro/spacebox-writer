package clickhouse

import (
	"database/sql"

	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	insertMultiSendMessage = `
		INSERT INTO spacebox.multisend_message 
		    (height, address_from, addresses_to, tx_hash, coins, msg_index)`
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

func (ch *Clickhouse) MultiSendMessage(val model.MultiSendMessage) (err error) {
	var (
		stmt *sql.Stmt
		tx   *sql.Tx

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
