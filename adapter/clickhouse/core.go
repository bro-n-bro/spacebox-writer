package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	jsoniter "github.com/json-iterator/go"

	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	insertTransactionQuery = `
		INSERT INTO spacebox.transaction 
		    (hash, height, success, messages, memo, signatures, signer_infos,
		     fee, signer, gas_wanted, gas_used, raw_log, logs)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
)

func (ch *Clickhouse) Transaction(val model.Transaction) (err error) {
	var (
		messages         = make([]interface{}, 0)
		feeBytes         []byte
		signerInfosBytes []byte
		messagesBytes    []byte
		signatures       = make(clickhouse.ArraySet, len(val.Signatures))
	)

	for i, s := range val.Signatures {
		signatures[i] = s
	}

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
	if _, err = ch.sql.Exec(
		insertTransactionQuery,
		val.Hash,
		val.Height,
		val.Success,
		string(messagesBytes),
		val.Memo,
		signatures,
		string(signerInfosBytes),
		string(feeBytes),
		val.Signer,
		val.GasWanted,
		val.GasUsed,
		val.RawLog,
		string(val.Logs),
	); err != nil {
		return err
	}

	return nil
}
