package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go/v2"
	jsoniter "github.com/json-iterator/go"

	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	insertTransactionQuery = `
		INSERT INTO spacebox.transaction (hash, height, success, messages, memo, signatures, signer_infos,
		                                  fee, signer, gas_wanted, gas_used, raw_log, logs)`
)

// Transaction is a method for saving transaction data to clickhouse
func (ch *Clickhouse) Transaction(vals []model.Transaction) (err error) {
	tx, err := ch.sql.Begin()
	if err != nil {
		return err
	}

	var (
		feeStr, signerInfosStr, messagesStr string
		signatures                          clickhouse.ArraySet
		messages                            []interface{}
	)

	stmt, err := tx.Prepare(insertTransactionQuery)
	if err != nil {
		return err
	}

	defer func() { _ = stmt.Close() }()

	for _, val := range vals {
		signatures = make(clickhouse.ArraySet, len(val.Signatures))

		for i, s := range val.Signatures {
			signatures[i] = s
		}

		if feeStr, err = jsoniter.MarshalToString(val.Fee); err != nil {
			return err
		}
		if signerInfosStr, err = jsoniter.MarshalToString(val.SignerInfos); err != nil {
			return err
		}

		messages = make([]interface{}, 0, len(val.Messages))
		for _, msg := range val.Messages {
			var tmp interface{}
			if err = jsoniter.Unmarshal(msg, &tmp); err != nil {
				return err
			}
			messages = append(messages, tmp)
		}

		if messagesStr, err = jsoniter.MarshalToString(messages); err != nil {
			return err
		}

		if _, err = stmt.Exec(
			val.Hash,
			val.Height,
			val.Success,
			messagesStr,
			val.Memo,
			signatures,
			signerInfosStr,
			feeStr,
			val.Signer,
			val.GasWanted,
			val.GasUsed,
			val.RawLog,
			string(val.Logs),
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}
