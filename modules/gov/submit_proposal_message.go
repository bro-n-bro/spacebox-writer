package gov

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// SubmitProposalMessageHandler is a handler for submit proposal message event
func SubmitProposalMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.SubmitProposalMessage](msgs)
	if err != nil {
		return err
	}

	return ch.SubmitProposalMessage(vals)
}
