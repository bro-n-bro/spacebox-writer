package distribution

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// CommunityPoolHandler is a handler for community pool event
func CommunityPoolHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.CommunityPool](msgs)
	if err != nil {
		return err
	}
	return ch.CommunityPool(vals)
}
