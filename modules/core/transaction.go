package core

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// TransactionHandler is a handler for transaction event
func TransactionHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.Transaction](msgs)
	if err != nil {
		return err
	}

	return ch.Transaction(vals)
}
