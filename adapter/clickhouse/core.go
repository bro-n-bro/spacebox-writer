package clickhouse

import (
	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/hexy-dev/spacebox-writer/adapter/clickhouse/models"
	"github.com/hexy-dev/spacebox/broker/model"
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
	var (
		involvedAccountsAddressesBytes []byte
	)

	if involvedAccountsAddressesBytes, err = jsoniter.Marshal(val.InvolvedAccountsAddresses); err != nil {
		return err
	}

	if err = ch.gorm.Table("message").Create(storageModel.Message{
		TransactionHash:           val.TransactionHash,
		Index:                     val.Index,
		Type:                      val.Type,
		Value:                     string(val.Value),
		InvolvedAccountsAddresses: string(involvedAccountsAddressesBytes),
		Signer:                    val.Signer,
		MsgIndex:                  val.MsgIndex,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (ch *Clickhouse) Transaction(val model.Transaction) (err error) {
	var (
		messages         = make([]interface{}, 0)
		feeBytes         []byte
		signaturesBytes  []byte
		signerInfosBytes []byte
		messagesBytes    []byte
	)

	if feeBytes, err = jsoniter.Marshal(val.Fee); err != nil {
		return err
	}

	if signaturesBytes, err = jsoniter.Marshal(val.Signatures); err != nil {
		return err
	}

	if signerInfosBytes, err = jsoniter.Marshal(val.SignerInfos); err != nil {
		return err
	}

	for _, msg := range val.Messages {
		var tmp interface{}
		if err = jsoniter.Unmarshal(msg, &tmp); err == nil {
			return err
		}
		messages = append(messages, tmp)
	}

	if messagesBytes, err = jsoniter.Marshal(messages); err != nil {
		return err
	}

	if err = ch.gorm.Table("transaction").Create(storageModel.Transaction{
		Messages:    string(messagesBytes),
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

func (ch *Clickhouse) LatestBlockHeight() (int64, error) {
	var lastHeight int64

	if err := ch.gorm.Select("height").Table("block").Order("height DESK").
		Limit(1).Scan(&lastHeight).Error; err != nil {
		return 0, err
	}

	return lastHeight, nil
}
