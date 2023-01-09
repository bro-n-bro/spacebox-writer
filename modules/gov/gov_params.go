package gov

import (
	"context"

	"spacebox-writer/adapter/clickhouse"
	storageModel "spacebox-writer/adapter/clickhouse/models"

	"github.com/hexy-dev/spacebox/broker/model"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func GovParamsHandler(ctx context.Context, msg []byte, ch *clickhouse.Clickhouse) error {
	val := model.GovParams{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	tallyParamsBytes, err := jsoniter.Marshal(val.TallyParams)
	if err != nil {
		return err
	}

	votingParamsBytes, err := jsoniter.Marshal(val.VotingParams)
	if err != nil {
		return err
	}

	depositParamsBytes, err := jsoniter.Marshal(val.DepositParams)
	if err != nil {
		return err
	}

	if err = ch.GetGormDB(ctx).Table("gov_params").Create(storageModel.GovParams{
		DepositParams: string(depositParamsBytes),
		VotingParams:  string(votingParamsBytes),
		TallyParams:   string(tallyParamsBytes),
		Height:        val.Height,
	}).Error; err != nil {
		return err
	}

	return nil
}
