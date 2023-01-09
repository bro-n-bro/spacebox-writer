package core

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func TransactionHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.Transaction{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	feeBytes, err := jsoniter.Marshal(val.Fee)
	if err != nil {
		return err
	}

	signaturesBytes, err := jsoniter.Marshal(val.Signatures)
	if err != nil {
		return err
	}

	signerInfosBytes, err := jsoniter.Marshal(val.SignerInfos)
	if err != nil {
		return err
	}

	var msgs = make([]interface{}, 0)
	for _, m := range val.Messages {
		var intefs interface{}
		if err = jsoniter.Unmarshal(m, &intefs); err != nil {
			return err
		}
		msgs = append(msgs, intefs)
	}

	messages, err := jsoniter.Marshal(msgs)
	if err != nil {
		return err
	}

	if err = ch.GetGormDB(ctx).Table("transaction").Create(storageModel.Transaction{
		Messages:    string(messages),
		Logs:        string(val.Logs),
		SignerInfos: string(signerInfosBytes),
		Signatures:  string(signaturesBytes),
		Fee:         string(feeBytes),
		Hash:        val.Hash,
		Height:      val.Height,
		Success:     val.Success,
		Memo:        val.Memo,
		Signer:      val.Signer,
		GasWanted:   val.GasWanted,
		GasUsed:     val.GasUsed,
		RawLog:      val.RawLog,
	}).Error; err != nil {
		return err
	}

	return nil
}
