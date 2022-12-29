package gov

import (
	"context"
	"github.com/hexy-dev/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"
)

func ProposalTallyResultHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.ProposalTallyResult{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	val2 := storageModel.ProposalTallyResult{
		ProposalID: int64(val.ProposalID),
		Yes:        val.Yes,
		No:         val.No,
		Abstain:    val.Abstain,
		NoWithVeto: val.NoWithVeto,
		Height:     val.Height,
	}

	if err := ch.GetGormDB(ctx).Table("proposal_tally_result").Create(val2).Error; err != nil {
		return err
	}

	return nil
}
