package clickhouse

import (
	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableHandleValidatorSignature = "handle_validator_signature"
)

// HandleValidatorSignature is a method for saving handle validator signature to clickhouse
func (ch *Clickhouse) HandleValidatorSignature(vals []model.HandleValidatorSignature) (err error) {
	var (
		burned string
	)

	batch := make([]storageModel.HandleValidatorSignature, len(vals))
	for i, val := range vals {
		if burned, err = jsoniter.MarshalToString(val.Burned); err != nil {
			return err
		}

		batch[i] = storageModel.HandleValidatorSignature{
			Address: val.Address,
			Power:   val.Power,
			Reason:  val.Reason,
			Jailed:  val.Jailed,
			Burned:  burned,
			Height:  val.Height,
		}
	}

	return ch.gorm.Table(tableHandleValidatorSignature).CreateInBatches(batch, len(batch)).Error
}
