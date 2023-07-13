package distribution

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// DistributionCommissionHandler is a handler for distribution commission event
func DistributionCommissionHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.DistributionCommission](msgs)
	if err != nil {
		return err
	}
	return ch.DistributionCommission(vals)
}
