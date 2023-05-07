package authz

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// GrantMessageHandler is a handler for grant message event
func GrantMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.GrantMessage](msgs)
	if err != nil {
		return err
	}

	return ch.GrantMessage(vals)
}
