package staking

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// DelegationMessageHandler is a handler for delegation message event
func DelegationMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.DelegationMessage](msgs)
	if err != nil {
		return err
	}
	return ch.DelegationMessage(vals)
}
