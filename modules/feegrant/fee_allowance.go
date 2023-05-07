package feegrant

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// FeeAllowanceHandler is a handler for fee allowance event
func FeeAllowanceHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.FeeAllowance](msgs)
	if err != nil {
		return err
	}

	return ch.FeeAllowance(vals)
}
