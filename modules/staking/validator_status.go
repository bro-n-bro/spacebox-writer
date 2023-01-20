package staking

import (
	"context"

	jsoniter "github.com/json-iterator/go"

	"github.com/bro-n-bro/spacebox-writer/internal/rep"
	"github.com/bro-n-bro/spacebox/broker/model"
)

func ValidatorStatusHandler(ctx context.Context, msg []byte, ch rep.Storage) error {
	val := model.ValidatorStatus{}
	if err := jsoniter.Unmarshal(msg, &val); err != nil {
		return err
	}

	return ch.ValidatorStatus(val)
}
