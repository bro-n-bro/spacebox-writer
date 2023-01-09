package gov

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
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

	if err = ch.GetGormDB(ctx).Table("proposal_deposit_message").Create(storageModel.ProposalDepositMessage{
		DepositorAddress: val.DepositorAddress,
		Coins:            string(coinsBytes),
		TxHash:           val.TxHash,
		ProposalID:       int64(val.ProposalID),
		Height:           val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}
