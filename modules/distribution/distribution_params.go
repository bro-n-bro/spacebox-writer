package distribution

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// DistributionParamsHandler is a handler for distribution params event
func DistributionParamsHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.DistributionParams](msgs)
	if err != nil {
		return err
	}
	return ch.DistributionParams(vals)
}
