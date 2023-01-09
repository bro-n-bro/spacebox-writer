package gov

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func ProposalTallyResultHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {

	val := model.ProposalTallyResult{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	if err := ch.GetGormDB(ctx).Table("proposal_tally_result").Create(storageModel.ProposalTallyResult{
		ProposalID: int64(val.ProposalID),
		Yes:        val.Yes,
		No:         val.No,
		Abstain:    val.Abstain,
		NoWithVeto: val.NoWithVeto,
		Height:     val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}
