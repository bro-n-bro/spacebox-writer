package clickhouse

import (
	"database/sql"

	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableGrantMessage = "grant_message"
	tableExecMessage  = "exec_message"
)

// GrantMessage is a method for saving grant message data to clickhouse
func (ch *Clickhouse) GrantMessage(vals []model.GrantMessage) (err error) {
	batch := make([]storageModel.GrantMessage, len(vals))
	for i, val := range vals {
		batch[i] = storageModel.GrantMessage{
			Grantee:  val.Grantee,
			MsgType:  val.MsgType,
			Height:   val.Height,
			TxHash:   val.TxHash,
			MsgIndex: val.MsgIndex,
			Expiration: sql.NullTime{
				Time:  val.Expiration,
				Valid: !val.Expiration.IsZero(),
			},
		}
	}

	return ch.gorm.Table(tableGrantMessage).CreateInBatches(batch, len(batch)).Error
}

// ExecMessage is a method for saving exec message data to clickhouse
func (ch *Clickhouse) ExecMessage(vals []model.ExecMessage) (err error) {
	var (
		messagesStr string
		messages    []interface{}
	)

	batch := make([]storageModel.ExecMessage, len(vals))
	for i, val := range vals {
		messages = make([]interface{}, 0, len(val.Msgs))
		for _, msg := range val.Msgs {
			var tmp interface{}
			if err = jsoniter.Unmarshal(msg, &tmp); err != nil {
				return err
			}
			messages = append(messages, tmp)
		}

		if messagesStr, err = jsoniter.MarshalToString(messages); err != nil {
			return err
		}

		batch[i] = storageModel.ExecMessage{
			Grantee:  val.Grantee,
			Msgs:     messagesStr,
			Height:   val.Height,
			TxHash:   val.TxHash,
			MsgIndex: val.MsgIndex,
		}
	}

	return ch.gorm.Table(tableExecMessage).CreateInBatches(batch, len(batch)).Error
}
