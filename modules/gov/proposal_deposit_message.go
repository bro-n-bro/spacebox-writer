package gov

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// ProposalDepositMessageHandler is a handler for proposal deposit message event
func ProposalDepositMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) error {
	vals, err := utils.ConvertMessages[model.ProposalDepositMessage](msgs)
	if err != nil {
		return err
	}

	return ch.ProposalDepositMessage(vals)
}
