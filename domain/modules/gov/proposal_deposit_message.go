package gov

import (
	"context"
	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
)

func ProposalDepositMessageHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.ProposalDepositMessage{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	coinsBytes, err := jsoniter.Marshal(val.Coins)
	if err != nil {
		return err
	}

	val2 := storageModel.ProposalDepositMessage{
		DepositorAddress: val.DepositorAddress,
		Coins:            string(coinsBytes),
		TxHash:           val.TxHash,
		ProposalID:       int64(val.ProposalID),
		Height:           val.Height,
	}

	if err = ch.GetGormDB(ctx).Table("proposal_deposit_message").Create(val2).Error; err != nil {
		return err
	}

	return nil
}
