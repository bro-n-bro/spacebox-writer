package gov

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// ProposalHandler is a handler for proposal event
func ProposalHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.Proposal{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	return ch.Proposal(val)
}
