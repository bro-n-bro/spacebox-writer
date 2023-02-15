package staking

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// UnbondingDelegationMessageHandler is a handler for unbonding delegation message event
func UnbondingDelegationMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.UnbondingDelegationMessage](msgs)
	if err != nil {
		return err
	}
	return ch.UnbondingDelegationMessage(vals)
}
