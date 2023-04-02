package bank

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// AccountBalanceHandler is a handler for account balance event
func AccountBalanceHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.AccountBalance](msgs)
	if err != nil {
		return err
	}

	return ch.AccountBalance(vals)
}
