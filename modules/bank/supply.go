package bank

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// SupplyHandler is a handler for supply event
func SupplyHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.Supply](msgs)
	if err != nil {
		return err
	}

	return ch.Supply(vals)
}
