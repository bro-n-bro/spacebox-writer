package staking

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// RedelegationHandle2r is a handler for redelegation event
func RedelegationHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.Redelegation](msgs)
	if err != nil {
		return err
	}
	return ch.Redelegation(vals)
}
