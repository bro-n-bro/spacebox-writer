package authz

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// AuthzGrantHandler is a handler for authz grant event
func AuthzGrantHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.AuthzGrant](msgs)
	if err != nil {
		return err
	}

	return ch.AuthzGrant(vals)
}
