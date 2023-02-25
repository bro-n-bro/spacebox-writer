package staking

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// StakingParamsHandler is a handler for staking params event
func StakingParamsHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.StakingParams](msgs)
	if err != nil {
		return err
	}
	return ch.StakingParams(vals)
}
