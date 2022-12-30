package gov

import (
	"context"
	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
)

func ProposalVoteMessageHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.ProposalVoteMessage{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	val2 := storageModel.ProposalVoteMessage{
		ProposalID:   int64(val.ProposalID),
		VoterAddress: val.VoterAddress,
		Option:       val.Option,
		Height:       val.Height,
	}

	if err := ch.GetGormDB(ctx).Table("proposal_vote_message").Create(val2).Error; err != nil {
		return err
	}

	return nil
}
