package authz

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// ExecMessageHandler is a handler for account balance event
func ExecMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.ExecMessage](msgs)
	if err != nil {
		return err
	}

	return ch.ExecMessage(vals)
}
