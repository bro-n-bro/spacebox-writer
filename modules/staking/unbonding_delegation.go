package staking

import (
	"context"

	jsoniter "github.com/json-iterator/go"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox/broker/model"
)

// UnbondingDelegationHandler is a handler for unbonding delegation event
func UnbondingDelegationHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.UnbondingDelegation{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	return ch.UnbondingDelegation(val)
}
