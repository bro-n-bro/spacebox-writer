package clickhouse

import (
	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableExecMessage = "exec_message"
)

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
