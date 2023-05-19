package staking

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// CreateValidatorMessageHandler is a handler for create validator message event
func CreateValidatorMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.CreateValidatorMessage](msgs)
	if err != nil {
		return err
	}
	return ch.CreateValidatorMessage(vals)
}
