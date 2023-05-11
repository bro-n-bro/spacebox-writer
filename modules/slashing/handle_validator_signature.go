package slashing

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// HandleValidatorSignatureHandler is a handler for handle validator signature event
func HandleValidatorSignatureHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.HandleValidatorSignature](msgs)
	if err != nil {
		return err
	}

	return ch.HandleValidatorSignature(vals)
}
