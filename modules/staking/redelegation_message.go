package staking

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// RedelegationMessageHandler is a handler for redelegation message event
func RedelegationMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.RedelegationMessage](msgs)
	if err != nil {
		return err
	}
	return ch.RedelegationMessage(vals)
}
