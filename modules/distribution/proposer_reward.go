package distribution

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// ProposerRewardHandler is a handler for proposer reward event
func ProposerRewardHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.ProposerReward](msgs)
	if err != nil {
		return err
	}
	return ch.ProposerReward(vals)
}
