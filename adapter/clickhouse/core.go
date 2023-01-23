package clickhouse

import (
	"context"
	"strings"
	"unsafe"

	"github.com/ClickHouse/clickhouse-go/v2"
	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

var (
	insertMessageQuery = `
INSERT INTO spacebox.message (transaction_hash, msg_index, type, signer, value, involved_accounts_addresses)
VALUES (?, ?, ?, ?, ?, ?);`

	insertTransactionQuery = `INSERT INTO spacebox.transaction (hash, height, success, messages, memo, signatures, 
                                   signer_infos, fee, signer, gas_wanted, gas_used, raw_log, logs, code)`

	insertTransactionQuery2 = `
INSERT INTO spacebox.transaction (hash, height, success, messages, memo, signatures, signer_infos, fee, signer, 
                                  gas_wanted, gas_used, raw_log, logs, code)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
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
	accountAddresses := make(clickhouse.ArraySet, len(val.InvolvedAccountsAddresses))
	for i, addr := range val.InvolvedAccountsAddresses {
		accountAddresses[i] = addr
	}

	if _, err := ch.sql.Exec(insertMessageQuery, val.TransactionHash, val.MsgIndex, val.Type, val.Signer, string(val.Value),
		accountAddresses); err != nil {

		return err
	}

	return nil
}

func (ch *Clickhouse) Transaction(val model.Transaction) (err error) {
	return ch.Transaction3(val)
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

// TODO: for tests
func (ch *Clickhouse) Transaction2(val model.Transaction) (err error) {
	var (
		messages         = make([]interface{}, len(val.Messages))
		signatures       = make(clickhouse.ArraySet, len(val.Signatures))
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

	for i, msg := range val.Messages {
		var tmp interface{}
		if err = jsoniter.Unmarshal(msg, &tmp); err != nil {
			return err
		}
		messages[i] = tmp
	}

	if messagesBytes, err = jsoniter.Marshal(messages); err != nil {
		return err
	}

	for i, sig := range val.Signatures {
		signatures[i] = sig
	}

	tx := storageModel.Transaction{
		RawLog:      val.RawLog,
		SignerInfos: string(signerInfosBytes),
		Logs:        string(val.Logs),
		Messages:    string(messagesBytes),
		Memo:        val.Memo,
		Signer:      val.Signer,
		Hash:        val.Hash,
		Fee:         string(feeBytes),
		Signatures:  signatures,
		GasUsed:     val.GasUsed,
		GasWanted:   val.GasWanted,
		Height:      val.Height,
		Success:     val.Success,
	}

	txBytes, err := jsoniter.Marshal(tx)
	if err != nil {
		return err
	}

	ch.log.Warn().Int("bytes length", len(txBytes)).Msgf("struct size: %v", unsafe.Sizeof(tx))

	toInsertMessages := strings.ReplaceAll(string(messagesBytes), "@", "")
	if err := ch.driverConn.Exec(
		context.Background(),
		insertTransactionQuery2,
		val.Hash,
		val.Height,
		val.Success,
		toInsertMessages,
		val.Memo,
		signatures,
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

	return nil
}

// TODO: for tests
func (ch *Clickhouse) Transaction3(val model.Transaction) (err error) {
	var (
		messages         = make([]interface{}, len(val.Messages))
		signatures       = make(clickhouse.ArraySet, len(val.Signatures))
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

	for i, msg := range val.Messages {
		var tmp interface{}
		if err = jsoniter.Unmarshal(msg, &tmp); err != nil {
			return err
		}
		messages[i] = tmp
	}

	if messagesBytes, err = jsoniter.Marshal(messages); err != nil {
		return err
	}

	for i, sig := range val.Signatures {
		signatures[i] = sig
	}

	tx := storageModel.Transaction{
		RawLog:      val.RawLog,
		SignerInfos: string(signerInfosBytes),
		Logs:        string(val.Logs),
		Messages:    string(messagesBytes),
		Memo:        val.Memo,
		Signer:      val.Signer,
		Hash:        val.Hash,
		Fee:         string(feeBytes),
		Signatures:  signatures,
		GasUsed:     val.GasUsed,
		GasWanted:   val.GasWanted,
		Height:      val.Height,
		Success:     val.Success,
	}

	txBytes, err := jsoniter.Marshal(tx)
	if err != nil {
		return err
	}

	ctx := context.Background()
	batch, err := ch.driverConn.PrepareBatch(ctx, insertTransactionQuery)
	if err != nil {
		return err
	}

	ch.log.Warn().Int("bytes length", len(txBytes)).Msgf("struct size: %v", unsafe.Sizeof(tx))
	toInsertMessages := strings.ReplaceAll(string(messagesBytes), "@", "")

	if err := batch.Append(
		val.Hash,
		val.Height,
		val.Success,
		toInsertMessages,
		val.Memo,
		signatures,
		string(signerInfosBytes),
		string(feeBytes),
		val.Signer,
		val.GasWanted,
		val.GasUsed,
		val.RawLog,
		string(val.Logs),
		uint32(0)); err != nil {

		return err
	}

	if err := batch.Send(); err != nil {
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
