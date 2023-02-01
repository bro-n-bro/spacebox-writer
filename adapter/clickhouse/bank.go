package clickhouse

import (
	"database/sql"

	"github.com/jinzhu/copier"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableMultiSendMessage = "multisend_message"
	tableAccountBalance   = "account_balance"
	tableSendMessage      = "send_message"
	tableSupply           = "supply"

	insertMultiSendMessage = `
		INSERT INTO spacebox.multisend_message 
		    (height, address_from, addresses_to, tx_hash, coins, msg_index)`
)

func (ch *Clickhouse) AccountBalance(val model.AccountBalance) (err error) {
	var (
		coinsBytes     []byte
		updates        storageModel.AccountBalance
		prevValStorage storageModel.AccountBalance
		valStorage     storageModel.AccountBalance
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	valStorage = storageModel.AccountBalance{
		Coins:   string(coinsBytes),
		Address: val.Address,
		Height:  val.Height,
	}

	if err = ch.gorm.Table(tableAccountBalance).
		Where("address = ?", valStorage.Address).
		First(&prevValStorage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = ch.gorm.Table(tableAccountBalance).Create(valStorage).Error; err != nil {
				return errors.Wrap(err, "error of create")
			}
			return nil
		}
		return errors.Wrap(err, "error of database: %w")
	}

	if valStorage.Height > prevValStorage.Height {
		if err = copier.Copy(&valStorage, &updates); err != nil {
			return errors.Wrap(err, "error of prepare update")
		}
		if err = ch.gorm.Table(tableAccountBalance).Where("address = ?", valStorage.Height).Updates(&updates).Error; err != nil {
			return errors.Wrap(err, "error of update")
		}
	}

	return nil
}

func (ch *Clickhouse) MultiSendMessage(val model.MultiSendMessage) (err error) {
	var (
		stmt *sql.Stmt
		tx   *sql.Tx

		coinsBytes []byte
		exists     bool
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	if exists, err = ch.ExistsTx(tableMultiSendMessage, val.TxHash, val.MsgIndex); err != nil {
		return err
	}

	if !exists {
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
	}
	return tx.Commit()
}

func (ch *Clickhouse) SendMessage(val model.SendMessage) (err error) {
	var (
		coinsBytes []byte
		exists     bool
	)

	if coinsBytes, err = jsoniter.Marshal(val.Coins); err != nil {
		return err
	}

	if exists, err = ch.ExistsTx(tableSendMessage, val.TxHash, val.MsgIndex); err != nil {
		return err
	}
	if !exists {
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

	if err = ch.gorm.Table(tableSupply).Create(storageModel.Supply{
		Coins:  string(coinsBytes),
		Height: val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}
