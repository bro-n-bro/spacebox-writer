package clickhouse

import (
	"database/sql"

	jsoniter "github.com/json-iterator/go"

	storageModel "github.com/bro-n-bro/spacebox-writer/adapter/clickhouse/models"
	"github.com/bro-n-bro/spacebox/broker/model"
)

const (
	tableFeeAllowance          = "fee_allowance"
	tableGrantAllowanceMessage = "grant_allowance_message"
)

// FeeAllowance is a method for saving fee allowance data to clickhouse
func (ch *Clickhouse) FeeAllowance(vals []model.FeeAllowance) (err error) {
	var (
		allowance string
	)

	batch := make([]storageModel.FeeAllowance, len(vals))
	for i, val := range vals {
		if allowance, err = jsoniter.MarshalToString(val.Allowance); err != nil {
			return err
		}

		batch[i] = storageModel.FeeAllowance{
			Grantee: val.Grantee,
			Granter: val.Granter,
			Height:  val.Height,
			Expiration: sql.NullTime{
				Time:  val.Expiration,
				Valid: !val.Expiration.IsZero(),
			},
			Allowance: allowance,
			IsActive:  val.IsActive,
		}
	}

	return ch.gorm.Table(tableFeeAllowance).CreateInBatches(batch, len(batch)).Error
}

// GrantAllowanceMessage is a method for saving account balance data to clickhouse
func (ch *Clickhouse) GrantAllowanceMessage(vals []model.GrantAllowanceMessage) (err error) {
	var (
		allowance string
	)

	batch := make([]storageModel.GrantAllowanceMessage, len(vals))
	for i, val := range vals {
		if allowance, err = jsoniter.MarshalToString(val.Allowance); err != nil {
			return err
		}

		batch[i] = storageModel.GrantAllowanceMessage{
			Granter:   val.Granter,
			Grantee:   val.Grantee,
			Allowance: allowance,
			TxHash:    val.TxHash,
			Height:    val.Height,
			MsgIndex:  val.MsgIndex,
			Expiration: sql.NullTime{
				Time:  val.Expiration,
				Valid: !val.Expiration.IsZero(),
			},
		}
	}

	return ch.gorm.Table(tableGrantAllowanceMessage).CreateInBatches(batch, len(batch)).Error
}
