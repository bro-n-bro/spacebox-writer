package feegrant

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// GrantAllowanceMessageHandler is a handler for grant allowance message event
func GrantAllowanceMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.GrantAllowanceMessage](msgs)
	if err != nil {
		return err
	}
	return ch.GrantAllowanceMessage(vals)
}
