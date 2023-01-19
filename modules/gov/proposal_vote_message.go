package gov

import (
	"context"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"

	"github.com/hexy-dev/spacebox-writer/internal/rep"
	"github.com/hexy-dev/spacebox/broker/model"
)

func ProposalVoteMessageHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.ProposalVoteMessage{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	return ch.ProposalVoteMessage(val)
}
