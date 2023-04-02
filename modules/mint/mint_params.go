package mint

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// MintParamsHandler is a handler for mint params event
func MintParamsHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.MintParams](msgs)
	if err != nil {
		return err
	}

	return ch.MintParams(vals)
}
