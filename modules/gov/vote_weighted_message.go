package gov

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox-writer/modules/utils"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// VoteWeightedMessageHandler is a handler for vote weighted message event
func VoteWeightedMessageHandler(ctx context.Context, msgs [][]byte, ch rep.Storage) (err error) {
	var vals []model.VoteWeightedMessage

	if vals, err = utils.ConvertMessages[model.VoteWeightedMessage](msgs); err != nil {
		return err
	}

	return ch.VoteWeightedMessage(vals)
}
