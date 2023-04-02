package staking

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// UnbondingDelegationHandler is a handler for unbonding delegation event
func UnbondingDelegationHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.UnbondingDelegation](msgs)
	if err != nil {
		return err
	}
	return ch.UnbondingDelegation(vals)
}
