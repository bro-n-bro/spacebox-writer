package gov

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func ProposalDepositHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.ProposalDeposit{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	coinsBytes, err := jsoniter.Marshal(val.Coins)
	if err != nil {
		return err
	}

	val2 := storageModel.ProposalDeposit{
		Coins:  string(coinsBytes),
		Height: val.Height,
	}

	if err = ch.GetGormDB(ctx).Table("proposal_deposit").Create(val2).Error; err != nil {
		return err
	}

	return nil
}
