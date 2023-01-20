package staking

import (
	"context"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox/broker/model"
	jsoniter "github.com/json-iterator/go"
)

func UnbondingDelegationMessageHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.UnbondingDelegationMessage{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	return ch.UnbondingDelegationMessage(val)
}
