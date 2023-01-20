package gov

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

func ProposalVoteMessageHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.ProposalVoteMessage{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return errors.Wrap(err, "unmarshall error")
	}

	return ch.ProposalVoteMessage(val)
}
