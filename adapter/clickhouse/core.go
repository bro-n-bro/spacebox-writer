package clickhouse

import (
	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
)

var (
	insertMessageQuery = `INSERT INTO spacebox.message (transaction_hash, msg_index, type, signer, value, 
                              involved_accounts_addresses)`

	insertTransactionQuery = `INSERT INTO spacebox.transaction (hash, height, success, messages, memo, signatures, 
                                   signer_infos, fee, signer, gas_wanted, gas_used, raw_log, logs, code)`
)

func (ch *Clickhouse) Block(val model.Block) error {
	if err := ch.gorm.Table("block").Create(storageModel.Block{
		Height:          val.Height,
		Hash:            val.Hash,
		NumTXS:          val.TxNum,
		TotalGas:        int64(val.TotalGas),
		ProposerAddress: val.ProposerAddress,
		Timestamp:       val.Timestamp,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) Message(val model.Message) (err error) {
	tx, err := ch.sql.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(insertMessageQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		val.TransactionHash,
		val.MsgIndex,
		val.Type,
		string(val.Value),
		val.Signer,
		val.InvolvedAccountsAddresses,
	); err != nil {
		return err
	}

	return tx.Commit()
}

func (ch *Clickhouse) Transaction(val model.Transaction) (err error) {
	var (
		messages         = make([]interface{}, 0)
		feeBytes         []byte
		signerInfosBytes []byte
		messagesBytes    []byte
	)

	if feeBytes, err = jsoniter.Marshal(val.Fee); err != nil {
		return err
	}

	if signerInfosBytes, err = jsoniter.Marshal(val.SignerInfos); err != nil {
		return err
	}

	for _, msg := range val.Messages {
		var tmp interface{}
		if err = jsoniter.Unmarshal(msg, &tmp); err != nil {
			return err
		}
		messages = append(messages, tmp)
	}

	if messagesBytes, err = jsoniter.Marshal(messages); err != nil {
		return err
	}

	tx, err := ch.sql.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(insertTransactionQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(
		val.Hash,
		val.Height,
		val.Success,
		string(messagesBytes),
		val.Memo,
		val.Signatures,
		string(signerInfosBytes),
		string(feeBytes),
		val.Signer,
		val.GasWanted,
		val.GasUsed,
		val.RawLog,
		string(val.Logs),
		uint32(0),
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) LatestBlockHeight() (int64, error) {
	var lastHeight int64

	if err := ch.gorm.Select("height").Table("block").Order("height DESC").
		Limit(1).Scan(&lastHeight).Error; err != nil {
		return 0, err
	}

	return lastHeight, nil
}
