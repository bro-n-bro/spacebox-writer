package staking

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// EditValidatorMessageHandler is a handler for edit validator message event
func EditValidatorMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.EditValidatorMessage](msgs)
	if err != nil {
		return err
	}
	return ch.EditValidatorMessage(vals)
}
