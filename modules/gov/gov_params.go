package gov

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// GovParamsHandler is a handler for gov params event
func GovParamsHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.GovParams](msgs)
	if err != nil {
		return err
	}

	return ch.GovParams(vals)
}
