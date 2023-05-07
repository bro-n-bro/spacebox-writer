package distribution

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// DistributionRewardHandler is a handler for distribution reward event
func DistributionRewardHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	dr, err := utils.ConvertMessages[model.DistributionReward](msgs)
	if err != nil {
		return err
	}
	return ch.DistributionReward(dr)
}
