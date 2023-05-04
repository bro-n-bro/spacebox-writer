package distribution

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// WithdrawValidatorCommissionMessageHandler is a handler for distribution params event
func WithdrawValidatorCommissionMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.WithdrawValidatorCommissionMessage](msgs)
	if err != nil {
		return err
	}
	return ch.WithdrawValidatorCommissionMessage(vals)
}
